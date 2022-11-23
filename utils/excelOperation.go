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

func (e *Excel) SetCellValues(sheets []string, column rune, row int, values []interface{}) error {
	// Case yg hanya 1 value
	if !(len(sheets) > 1) && !(len(values) > 1) {
		axis := string(column) + strconv.Itoa(row)
		err := e.excelFile.SetCellValue(sheets[0], axis, values[0])
		if err != nil {
			return err
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
				return err
			}
			startColumn++
			row++
		}
		startColumn = initColVal
		row = initRowVal
	}
	return nil
}

func (e *Excel) SaveFile() error {
	err := e.excelFile.SaveAs(e.fileName, excelize.Options{Password: e.password})
	if err != nil {
		return err
	}
	return nil
}
