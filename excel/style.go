package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// styleConfig 设置Excel的样式，如对齐，列宽等
func styleConfig(f *excelize.File, sheetName string, rowNum, colNum int, colWidth []float64) error {
	endCol, err := excelize.ColumnNumberToName(colNum)
	if err != nil {
		return err
	}
	// 设置title背景颜色
	titleStyle, err := getTitleStyle(f)
	if err != nil {
		return err
	}
	err = f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s1", endCol), titleStyle)
	if err != nil {
		return err
	}

	// 设置除了标题以外的 style
	contentStyle, err := getContentStyle(f)
	if err != nil {
		return err
	}
	err = f.SetCellStyle(sheetName, "A2", fmt.Sprintf("%s%d", endCol, rowNum), contentStyle)
	if err != nil {
		return err
	}

	// 设置列宽
	for col, colWidth := range colWidth {
		colName, err := excelize.ColumnNumberToName(col + 1)
		if err != nil {
			return err
		}
		err = f.SetColWidth(sheetName, colName, colName, colWidth)
		if err != nil {
			return err
		}
	}

	return nil
}

func getTitleStyle(f *excelize.File) (styleID int, err error) {
	styleID, err = f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E0EBF5"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			WrapText:   true,
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	return
}
func getContentStyle(f *excelize.File) (int, error) {
	return f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical:        "center",
			Horizontal:      "center",
			JustifyLastLine: true,
			WrapText:        true,
		},
	})
}
