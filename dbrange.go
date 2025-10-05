// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssdb

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/api/sheets/v4"
)

type DBRange struct {
	symbolicRange string
	gridRange     *sheets.GridRange
	sheet         *Sheet
}

func (dbrange *DBRange) Extend(datum *DBRange) (res *DBRange) {
	foo := *dbrange
	res = &foo
	if datum.gridRange.StartColumnIndex < dbrange.gridRange.StartColumnIndex {
		res.gridRange.StartColumnIndex = datum.gridRange.StartColumnIndex
	}
	if datum.gridRange.StartRowIndex < dbrange.gridRange.StartRowIndex {
		res.gridRange.StartRowIndex = datum.gridRange.StartRowIndex
	}
	if datum.gridRange.EndColumnIndex > dbrange.gridRange.EndColumnIndex {
		res.gridRange.EndColumnIndex = datum.gridRange.EndColumnIndex
	}
	if datum.gridRange.EndRowIndex > dbrange.gridRange.EndRowIndex {
		res.gridRange.EndRowIndex = datum.gridRange.EndRowIndex
	}
	return //
}

func (dbrange *DBRange) NeedsGrowth(datum *DBRange) (growRows, growColumns int64) {
	if datum.gridRange.EndRowIndex > dbrange.gridRange.EndRowIndex {
		growRows = datum.gridRange.EndRowIndex - dbrange.gridRange.EndRowIndex
	}
	if datum.gridRange.EndColumnIndex > dbrange.gridRange.EndColumnIndex {
		growColumns = datum.gridRange.EndColumnIndex - dbrange.gridRange.EndColumnIndex
	}
	return //
}

func (dbrange *DBRange) String() (s string) {
	s += dbrange.sheet.Sheet.Properties.Title + "!"
	rowMin, colMin, rowMax, colMax :=
		dbrange.gridRange.StartRowIndex,
		dbrange.gridRange.StartColumnIndex,
		dbrange.gridRange.EndRowIndex,
		dbrange.gridRange.EndColumnIndex

	s += FormatNumericColumn(false, colMin)
	s += FormatNumericRow(false, rowMin)

	swrow := FormatNumericRow(true, rowMax)
	swcol := FormatNumericColumn(true, colMax)
	if len(swrow+swcol) > 0 {
		s += ":" + swcol + swrow
	}
	return //
}

func (cell *Cell) Range() (dbRange *DBRange) {
	if cell == nil {
		return nil
	}
	gridRange := &sheets.GridRange{
		SheetId:          cell.Sheet.GetID(),
		StartColumnIndex: cell.N,
		EndColumnIndex:   cell.N + 1,
		StartRowIndex:    cell.Row.N,
		EndRowIndex:      cell.Row.N + 1,
	}
	dbRange = &DBRange{
		gridRange: gridRange,
		sheet:     cell.Sheet,
	}
	dbRange.symbolicRange = dbRange.String()
	return //
}

func (ssdb *SSDB) NewDBRange(sheetname string, row, col, rows, cols int64) (dbRange *DBRange) {
	sheet := ssdb.SheetLookup(sheetname)
	if sheet == nil {
		return nil
	}
	gridRange := &sheets.GridRange{
		SheetId:          sheet.GetID(),
		StartColumnIndex: col,
		EndColumnIndex:   col + cols,
		StartRowIndex:    row,
		EndRowIndex:      row + rows,
	}
	dbRange = &DBRange{
		gridRange: gridRange,
		sheet:     sheet,
	}
	dbRange.symbolicRange = dbRange.String()
	return //
}

func (ssdb *SSDB) SheetLookup(sheetMatch string) (foundSheet *Sheet) {
	err := false
	ssdb.SheetIter(func(sheetname string, sheet *Sheet) {
		switch {
		case foundSheet != nil && sheetname == sheetMatch:
			err = true
		case sheetname == sheetMatch:
			foundSheet = sheet
		}
	})
	if err || foundSheet == nil {
		return nil
	}
	return //
}

func (ssdb *SSDB) NewDBRangeFromSymbolicRange(symbolicRange string) (dbRange *DBRange) {
	spl := strings.Split(symbolicRange, "!")
	if len(spl) < 2 {
		return nil
	}
	sheetMatch := spl[0]
	row0, col0 := ParseSymbolicRange(false, spl[1])
	row1, col1 := ParseSymbolicRange(true, spl[1])
	if row0 == -1 || col0 == -1 || row1 == -1 || col1 == -1 {
		return nil
	}
	gridRange := &sheets.GridRange{
		StartColumnIndex: col0,
		EndColumnIndex:   col1,
		StartRowIndex:    row0,
		EndRowIndex:      row1,
	}
	var foundSheet *Sheet
	err := false
	ssdb.SheetIter(func(sheetname string, sheet *Sheet) {
		switch {
		case foundSheet != nil && sheetname == sheetMatch:
			err = true
		case sheetname == sheetMatch:
			foundSheet = sheet
			gridRange.SheetId = sheet.GetID()
		}
	})
	if err || foundSheet == nil {
		return nil
	}
	dbRange = &DBRange{
		symbolicRange: symbolicRange,
		gridRange:     gridRange,
		sheet:         foundSheet,
	}
	return //
}

func ParseSymbolicRangeNW(cellName string) (row, col int64) {
	return ParseSymbolicRange(false, cellName)
}

func ParseSymbolicRangeSE(cellName string) (row, col int64) {
	return ParseSymbolicRange(true, cellName)
}

// corner: false UL, true SE
func ParseSymbolicRange(corner bool, cellName string) (row, col int64) {
	spl := strings.Split(cellName, ":")
	switch {
	case len(spl) == 0:
		return -1, -1
	case len(spl) == 1:
		return ParseSymbolicCell(corner, cellName)
	case !corner && len(spl) == 2:
		return ParseSymbolicCell(corner, spl[0])
	case corner && len(spl) == 2:
		return ParseSymbolicCell(corner, spl[1])
	}
	return //
}

// corner: false UL, true SE
func ParseSymbolicCell(corner bool, cellName string) (row, col int64) {
	alpha := ""
	numeric := ""
	for _, c := range cellName {
		switch {
		case (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') && len(numeric) == 0:
			alpha += string(c)
		case c >= '0' && c <= '9':
			numeric += string(c)
		default:
			return -1, -1
		}
	}
	switch {
	case !corner && len(numeric) == 0 && len(alpha) == 0:
		return 0, 0
	case !corner && len(numeric) == 0:
		return 0, ParseSymbolicColumn(alpha)
	case !corner && len(alpha) == 0:
		return ParseSymbolicRow(numeric), 0
	case corner && len(numeric) == 0 && len(alpha) == 0:
		return 9999, 9999
	case corner && len(numeric) == 0:
		return 9999, ParseSymbolicColumn(alpha) + 1
	case corner && len(alpha) == 0:
		return ParseSymbolicRow(numeric) + 1, 9999
	case corner:
		return ParseSymbolicRow(numeric) + 1, ParseSymbolicColumn(alpha) + 1
	default:
		return ParseSymbolicRow(numeric), ParseSymbolicColumn(alpha)
	}
}

// Accepts input A, B, C, ..., Z, AA, AB, ... AZ, ..., ZZ
func ParseSymbolicRow(rowNumStr string) (row int64) {
	if len(rowNumStr) == 0 || len(rowNumStr) > 5 {
		return -1
	}
	for _, c := range rowNumStr {
		switch {
		case c >= 48 && c < 58:
			row = row*10 + int64(c-48)
		default:
			return -1
		}
	}
	row--  // make zero based
	return //
}

// Accepts input A, B, C, ..., Z, AA, AB, ... AZ, ..., ZZ
func ParseSymbolicColumn(colName string) (col int64) {
	colName = strings.ToUpper(colName)
	switch {
	case len(colName) == 1:
		col = int64([]byte(colName)[0] - 65)
	case len(colName) == 2:
		col = int64(([]byte(colName)[0]-65)*26+[]byte(colName)[1]-65) + 26
	default:
		col = -1
	}
	return //
}

// Accepts input A, B, C, ..., Z, AA, AB, ... AZ, ..., ZZ
func FormatNumericRow(corner bool, rownum int64) (rowName string) {
	switch {
	case rownum < 0:
		panic("Out of range")
	case rownum < 9999:
		rowName = strconv.FormatInt(rownum, 10)
	case !corner && rownum == 9999:
		rowName = strconv.FormatInt(rownum, 10)
	case corner && rownum == 9999:
		rowName = ""
	default:
		rowName = "OUTOFRANGE"
		panic("colName out of range")
	}
	return //
}

// Accepts input A, B, C, ..., Z, AA, AB, ... AZ, ..., ZZ
func FormatNumericColumn(corner bool, colnum int64) (colName string) {
	switch {
	case colnum < 26:
		// colName = string('A' + colnum)
		colName = fmt.Sprintf("%c", colnum+65)
	case colnum < 26*26:
		// colName = string('A'+colnum/26) + string('A'+colnum%26)
		colName = fmt.Sprintf("%c%c", colnum/26+65, colnum%26+65)
	case corner && colnum == 9999:
		colName = ""

	// error conditions:
	case !corner && colnum == 9999:
		fallthrough
	default:
		colName = "OUTOFRANGE"
		panic("colName out of range")
	}
	return //
}
