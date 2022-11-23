package utils

import (
	"github.com/xuri/excelize/v2"
	"strconv"
)

type Excel struct {
	excelFile *excelize.File
	fileName  string
	password  string
}

type Option func(e *Excel)

func WithPassword(pass string) Option {
	return func(e *Excel) {
		e.password = pass
	}
}

func NewExcel(fileName string, opts ...Option) *Excel {
	e := &Excel{excelFile: excelize.NewFile(), fileName: fileName, password: ""}

	// Apply Option
	for _, opt := range opts {
		opt(e)
	}

	return e
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

func (e *Excel) ChangeSheetName(oldName, newName string) {
	e.excelFile.SetSheetName(oldName, newName)
}

func (e *Excel) SetCellValues(sheets []string, column rune, row int, values []interface{}) {
	// Case yg hanya 1 value
	if !(len(sheets) > 1) && !(len(values) > 1) {
		axis := string(column) + strconv.Itoa(row)
		err := e.excelFile.SetCellValue(sheets[0], axis, values[0])
		if err != nil {
			panic(err)
		}
	}

	startColumn := int(column)
	initRowVal := row
	initColVal := startColumn
	for _, sheet := range sheets {
		for _, value := range values {
			columnStr := string(rune(startColumn))
			rowStr := strconv.Itoa(row)
			axis := columnStr + rowStr
			err := e.excelFile.SetCellValue(sheet, axis, value)
			if err != nil {
				panic(err)
			}
			startColumn++
		}
		startColumn = initColVal
		row = initRowVal
	}
}

func (e *Excel) MakeStyle(style interface{}) int {
	styleInt, err := e.excelFile.NewStyle(style)

	if err != nil {
		panic(err)
	}
	return styleInt
}

func (e *Excel) SetColWidth(sheet, startCol, endCol string, width float64) {
	err := e.excelFile.SetColWidth(sheet, startCol, endCol, width)
	if err != nil {
		panic(err)
	}
}

func (e *Excel) SetCellStyle(sheet, hCell, vCell string, styleID int) {
	err := e.excelFile.SetCellStyle(sheet, hCell, vCell, styleID)
	if err != nil {
		panic(err)
	}
}

func (e *Excel) SaveFile() error {
	err := e.excelFile.SaveAs(e.fileName, excelize.Options{Password: e.password})
	if err != nil {
		return err
	}
	return nil
}
