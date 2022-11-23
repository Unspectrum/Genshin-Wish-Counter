package utils

import "github.com/xuri/excelize/v2"

type Excel struct {
	excelFile *excelize.File
}

func NewExcel() *Excel {
	return &Excel{excelFile: excelize.NewFile()}
}

func (e *Excel) ExcelFile() *excelize.File {
	return e.excelFile
}

func (e *Excel) GenerateSheets(sheets []string) int {
	for _, sheet := range sheets {
		e.ExcelFile().NewSheet(sheet)
	}
	return e.ExcelFile().GetSheetIndex(sheets[len(sheets)-1])
}
