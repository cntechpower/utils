package excel

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"

	"github.com/cntechpower/utils/oss"
)

const (
	maxCellWidth         = 50
	commonCellWidth      = 20
	defaultTagSeparator  = "|"
	defaultPictureHeight = 100
	defaultPictureWidth  = 100
	defaultSheetName     = "Sheet1"
	ExtraTypePicture     = "picture"
)

type Extra struct {
	Type   string
	Params []string
	Data   interface{}
}

// SheetData 写入Excel的sheet data
type SheetData struct {
	SheetName string // 可省略，省略就以Sheet1,Sheet2...命名
	// 需要是一个二维切片[][]interface{}或者[]Material{}自定义结构体切片
	// 特别说明，如果data传的是[]Struct的情况下，假设某个字段是State，并且定义了func (Struct) StateString() string，
	// 那么写入Excel的时候会自动调用该方法，代替State的原始值。
	// 注意是func (Struct) StateString() string不是func (*Struct) StateString() string
	// 主要用于像一些状态1，2导出的时候给出中文意思：1-启用；2-禁用
	Data  interface{}
	Extra map[string]*Extra
}

func parseExcelTag(tag string) map[string]interface{} {
	if tag == "" {
		return nil
	}
	items := strings.Split(tag, defaultTagSeparator)
	res := make(map[string]interface{})
	for _, item := range items {
		tagItemInfo := strings.Split(item, "=")
		if len(tagItemInfo) < 2 {
			continue
		}
		key := tagItemInfo[0]
		value := tagItemInfo[1]
		res[key] = value
	}
	return res
}

const (
	titleTagKeyEn = "title_en"
	titleTagKeyCN = "title"
)

func getTitleFromTagInfo(tagInfo map[string]interface{}) (title string, ok bool) {
	title, ok = tagInfo[titleTagKeyCN].(string)
	if !ok {
		return
	}

	titleEN, ok1 := tagInfo[titleTagKeyEn].(string)
	if ok1 {
		title = title + "\r\n" + titleEN
	}

	return
}

func getTitleFromStruct(t reflect.Type) (title []interface{}, selectedField []int) {
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("excel")
		if tag == "" {
			continue
		}
		tagInfo := parseExcelTag(tag)
		if name, ok := getTitleFromTagInfo(tagInfo); ok {
			title = append(title, name)
			selectedField = append(selectedField, i)
		}
	}
	return
}

func transformStruct(ctx context.Context, data interface{}) (interface{}, map[string]*Extra, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return data, nil, fmt.Errorf("excelData.SheetData 不是一个slice")
	}
	if v.Len() == 0 {
		return data, nil, fmt.Errorf("excelData.SheetData 为空")
	}
	vt := v.Index(0)
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}
	// 检查data[0]是否是slice
	if vt.Kind() == reflect.Slice {
		return data, nil, nil
	} else if vt.Kind() != reflect.Struct {
		return data, nil, errors.New("excelData.Data不是一个正确的数据，期望[][]interface{}或[]MyStruct{}")
	}
	// struct 类型
	// 将data.Data转为[][]interface{}
	sliceData := make([][]interface{}, 0, v.Len()+1)
	title, selected := getTitleFromStruct(vt.Type())
	sliceData = append(sliceData, title)
	extra := make(map[string]*Extra, 0)
	startRow := len(sliceData)
	for index := 0; index < v.Len(); index++ {
		sv := v.Index(index)
		if sv.Kind() == reflect.Ptr {
			sv = sv.Elem()
		}
		row := make([]interface{}, 0, len(selected))
		for i := 0; i < sv.NumField(); i++ {
			if lo.Contains(selected, i) {
				field := sv.Type().Field(i)
				methodName := fmt.Sprintf("%sString", field.Name)
				method := sv.MethodByName(methodName)
				if method.IsValid() {
					var res []reflect.Value
					if method.Type().NumIn() == 1 {
						res = method.Call([]reflect.Value{reflect.ValueOf(ctx)})
					} else {
						res = method.Call(nil)
					}

					if len(res) > 0 {
						row = append(row, res[0].Interface())
					}
				} else {
					if extraTag := field.Tag.Get("extra"); extraTag != "" {
						axis, err := excelize.CoordinatesToCellName(i+1, index+startRow+1)
						if err != nil {
							return nil, nil, err
						}

						extraArgs := strings.Split(extraTag, ",")
						extra[axis] = &Extra{
							Type:   extraArgs[0],
							Params: extraArgs[1:],
							Data:   sv.Field(i).Interface(),
						}

						row = append(row, "")
					} else {
						row = append(row, sv.Field(i).Interface())
					}
				}
			}
		}
		sliceData = append(sliceData, row)
	}
	return sliceData, extra, nil
}

func checkDataOK(ctx context.Context, excelData []SheetData) error {
	// 数据简单校验
	for i, data := range excelData {
		res, extra, err := transformStruct(ctx, data.Data)
		if err != nil {
			return err
		}
		excelData[i].Data = res
		excelData[i].Extra = extra
	}
	return nil
}

func storeColWidth(colsWidth []float64, colIndex int, value interface{}) {
	colLen := len(colsWidth)
	if colLen == 0 {
		panic("请传递一个长度等于列数量的slice进来")
	}
	if colIndex >= colLen || colIndex < 0 {
		return
	}
	var calculateWidth float64 = commonCellWidth
	switch value := value.(type) {
	case time.Time:
		calculateWidth = 19
	case string:
		calculateWidth = float64(len(value))
	}
	if calculateWidth > maxCellWidth {
		calculateWidth = maxCellWidth
	}
	if calculateWidth > colsWidth[colIndex] {
		colsWidth[colIndex] = calculateWidth
	}
}

// PixelToRowHeight 将像素高度转换为Excel行高宽(估算，并不会准确)
func PixelToRowHeight(pixelHeight int) float64 {
	rowHeight := float64(pixelHeight) / 1.20
	return rowHeight
}

// PixelToColWidth 将像素宽度转换为Excel列宽(估算，并不会准确)
func PixelToColWidth(pixelWidth int) float64 {
	colWidth := float64(pixelWidth-5) / 7
	return colWidth
}

func writeSheetData(ctx context.Context, f *excelize.File, sheetName string, data interface{}, extra map[string]*Extra) error {
	// 检查数据是否为slice
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return errors.New("data 不是一个slice")
	}
	if v.Len() == 0 {
		return errors.New("data 为空")
	}
	// 检查data[0]是否是slice
	if v.Index(0).Kind() != reflect.Slice {
		return errors.New("data 不是一个二维切片")
	}
	rowLen := v.Len()
	colLen := v.Index(0).Len()
	// 记录每一列最大文本长度
	colMaxWidth := make([]float64, colLen)
	for row := 0; row < v.Len(); row++ {
		rowData := v.Index(row)
		for col := 0; col < rowData.Len(); col++ {
			cellValue := rowData.Index(col).Interface()
			// 如果是日期，转换下格式
			if t, ok := cellValue.(time.Time); ok {
				cellValue = t.Format("2006-01-02 15:04:05")
			}
			// 设置长度
			storeColWidth(colMaxWidth, col, cellValue)
			axis, err := excelize.CoordinatesToCellName(col+1, row+1)
			if err != nil {
				return err
			}

			if ex, ok := extra[axis]; ok {
				customWidth, err := setCellExtra(ctx, f, sheetName, axis, ex)
				if err != nil {
					return err
				}

				colMaxWidth[col] = customWidth
			} else {
				err = f.SetCellValue(sheetName, axis, cellValue)
				if err != nil {
					return err
				}
			}
		}
	}
	// 设置样式
	err := styleConfig(f, sheetName, rowLen, colLen, colMaxWidth)
	if err != nil {
		return err
	}
	return nil
}

// WriteExcel 生成一个Excel文件，常规方式写入，建议少量数据导出使用
func WriteExcel(ctx context.Context, filename string, excelData ...SheetData) error {
	// 数据简单校验
	err := checkDataOK(ctx, excelData)
	if err != nil {
		return err
	}
	tmpChan := make(chan error, 1)
	go func() {
		f := excelize.NewFile()
		for index, data := range excelData {
			if data.SheetName == "" {
				data.SheetName = fmt.Sprintf("Sheet%d", index+1)
			}
			if index == 0 {
				_ = f.SetSheetName(defaultSheetName, data.SheetName)
			} else {
				_, _ = f.NewSheet(data.SheetName)
			}
			err = writeSheetData(ctx, f, data.SheetName, data.Data, data.Extra)
			if err != nil {
				tmpChan <- err
				return
			}
		}
		tmpChan <- f.SaveAs(filename)
	}()
	select {
	case <-ctx.Done():
		return errors.New("operation canceled")
	case err := <-tmpChan:
		return err
	}
}

func WriteExcelAndUploadOSS(ctx context.Context, filename string, excelData ...SheetData) (key string, err error) {
	err = WriteExcel(ctx, filename, excelData...)
	if err != nil {
		return
	}
	key, err = oss.Impl.Upload(ctx, filename)
	return
}
