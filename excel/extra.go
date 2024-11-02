package excel

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func fetchImageFile(ctx context.Context, urlAddr string) ([]byte, error) {
	// 使用HTTP客户端下载文件
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlAddr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	return bs, nil
}

func setCellExtra(ctx context.Context, f *excelize.File, sheetName string, axis string, extra *Extra) (widthChar float64, err error) {
	switch extra.Type {
	case ExtraTypePicture:
		uri, ok := extra.Data.(string)
		if !ok || uri == "" {
			return
		}

		var height = defaultPictureHeight
		var width = defaultPictureWidth
		if len(extra.Params) > 1 {
			heightValue, err := strconv.ParseInt(extra.Params[1], 10, 64)
			if err != nil {
				return 0, err
			}

			height = int(heightValue)
		}

		if len(extra.Params) > 2 {
			widthValue, err := strconv.ParseInt(extra.Params[2], 10, 64)
			if err != nil {
				return 0, err
			}

			width = int(widthValue)
		}

		widthChar = PixelToColWidth(width)
		heightPt := PixelToRowHeight(height)

		file, err := fetchImageFile(ctx, uri)
		if err != nil {
			return 0, err
		}

		u, err := url.Parse(uri)
		if err != nil {
			return 0, err
		}

		if err := excelSetPicture(f, sheetName, axis, heightPt, widthChar, u.Path, file); err != nil {
			return 0, err
		}
	}

	return widthChar, nil
}

func excelSetPicture(f *excelize.File, sheetName string, axis string, height, width float64, fileName string, file []byte) (err error) {
	var col, row int
	col, row, err = excelize.CellNameToCoordinates(axis)
	if err != nil {
		return err
	}

	if err := f.SetRowHeight(sheetName, row, height); err != nil {
		return err
	}

	// 由于图片的autofix是在设置的时候进行计算的，先设置一次宽度
	colName, err := excelize.ColumnNumberToName(col)
	if err != nil {
		return err
	}

	err = f.SetColWidth(sheetName, colName, colName, width)
	if err != nil {
		return err
	}

	if err := f.AddPictureFromBytes(sheetName, axis, &excelize.Picture{
		Extension: strings.ToLower(path.Ext(fileName)),
		File:      file,
		Format: &excelize.GraphicOptions{
			AutoFit: true,
		},
	}); err != nil {
		return err
	}

	return
}
