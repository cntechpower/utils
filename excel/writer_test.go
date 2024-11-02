package excel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cntechpower/utils/oss"
)

type TestStruct struct {
	Name   string    `excel:"title=姓名"`
	Age    int       `excel:"title=年龄"`
	Birth  time.Time `excel:"title=创建日期"`
	Status int       `excel:"title=状态"`
}

func TestWriteExcel(t *testing.T) {
	ctx := context.Background()
	filename := "test.xlsx"
	data := []TestStruct{
		{
			Name:   "张三",
			Age:    30,
			Birth:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
			Status: 1,
		},
		{
			Name:   "李四",
			Age:    25,
			Birth:  time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local),
			Status: 2,
		},
	}
	sheetData := []SheetData{
		{
			SheetName: "Sheet1",
			Data:      data,
		},
	}
	err := WriteExcel(ctx, filename, sheetData...)
	assert.NoError(t, err)
	storage := oss.NewMinio()
	key, err := storage.Upload(ctx, filename)
	assert.NoError(t, err)
	url, err := storage.GetURL(ctx, key)
	assert.NoError(t, err)
	t.Log(url)
}

func TestWriteExcel_Canceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	filename := "test.xlsx"
	data := []TestStruct{
		{
			Name:   "张三",
			Age:    30,
			Birth:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
			Status: 1,
		},
		{
			Name:   "李四",
			Age:    25,
			Birth:  time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local),
			Status: 2,
		},
	}
	sheetData := []SheetData{
		{
			SheetName: "Sheet1",
			Data:      data,
		},
	}
	err := WriteExcel(ctx, filename, sheetData...)
	assert.Error(t, err, "operation canceled")
}

func TestWriteExcel_InvalidData(t *testing.T) {
	ctx := context.Background()
	filename := "test.xlsx"
	data := []TestStruct{
		{
			Name:   "张三",
			Age:    30,
			Birth:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
			Status: 1,
		},
		{
			Name:   "李四",
			Age:    25,
			Birth:  time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local),
			Status: 2,
		},
	}
	sheetData := []SheetData{
		{
			SheetName: "Sheet1",
			Data:      data,
		},
	}
	err := WriteExcel(ctx, filename, sheetData...)
	assert.Error(t, err, "data 不是一个正确的数据，期望[][]interface{}或[]MyStruct{}")
}

func TestWriteExcel_EmptyData(t *testing.T) {
	ctx := context.Background()
	filename := "test.xlsx"
	data := []TestStruct{}
	sheetData := []SheetData{
		{
			SheetName: "Sheet1",
			Data:      data,
		},
	}
	err := WriteExcel(ctx, filename, sheetData...)
	assert.Error(t, err, "data 为空")
}

func TestWriteExcel_NotSliceData(t *testing.T) {
	ctx := context.Background()
	filename := "test.xlsx"
	data := TestStruct{
		Name:   "张三",
		Age:    30,
		Birth:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
		Status: 1,
	}
	sheetData := []SheetData{
		{
			SheetName: "Sheet1",
			Data:      data,
		},
	}
	err := WriteExcel(ctx, filename, sheetData...)
	assert.Error(t, err, "data 不是一个slice")
}

func (t TestStruct) StatusString() string {
	switch t.Status {
	case 1:
		return "启用"
	case 2:
		return "禁用"
	default:
		return "未知"
	}
}
