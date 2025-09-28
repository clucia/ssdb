// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssdb

import (
	"errors"
	"strings"

	"google.golang.org/api/sheets/v4"
)

func (sheet *Sheet) GetRange(dbrange *DBRange) (_data [][]any) {
	return sheet.CopyVals(dbrange)
}

func (sheet *Sheet) GetExtents() (res *DBRange) {
	res = &DBRange{
		gridRange: &sheets.GridRange{
			SheetId:          sheet.GetID(),
			StartRowIndex:    0,
			EndRowIndex:      0,
			StartColumnIndex: 0,
			EndColumnIndex:   0,
		},
	}
	sheet.RowIter(func(row *Row) {
		if res.gridRange.EndRowIndex < (row.N + 1) {
			res.gridRange.EndRowIndex = row.N + 1
		}
		row.CellIter(func(cell *Cell) {
			if res.gridRange.EndColumnIndex < (cell.N + 1) {
				res.gridRange.EndColumnIndex = cell.N + 1
			}
		})
	})
	return //
}

func (sheet *Sheet) GetRowN(N int64) *Row {
	if sheet == nil {
		return nil
	}
	return &Row{
		N:     N,
		DB:    sheet.DB,
		Sheet: sheet,
		Row:   sheet.Sheet.Data[0].RowData[N],
	}
}

var ErrDuplicateColumnKey = errors.New("duplicate column error")

func (row *Row) GetCellByName(name string) (res *Cell, err error) {
	hdrrow := row.Sheet.GetRowN(0)
	hdrrow.CellIter(func(cell *Cell) {
		switch {
		case err != nil:
			// do nothing after first error
		case res != nil && cell.GetString() == name:
			err = ErrDuplicateColumnKey
		case cell.GetString() == name:
			res = row.GetCellN(cell.N)
		}
	})
	if err != nil {
		return nil, err
	}
	return //
}

func (row *Row) GetCellN(N int64) *Cell {
	if row == nil || row.Row == nil {
		return nil
	}
	if len(row.Row.Values) <= int(N) {
		return nil
	}
	return &Cell{
		N:     N,
		DB:    row.DB,
		Sheet: row.Sheet,
		Row:   row,
		Cell:  row.Row.Values[N],
	}
}
func (row *Row) IsBlank() (res bool) {
	switch {
	case row == nil:
		return true
	case row.Row == nil:
		return true
	case row.Row.Values == nil:
		return true
	case len(row.Row.Values) == 0:
		return true
	}
	res = true
	row.CellIter(func(cell *Cell) {
		if cell.GetString() != "" {
			res = false
		}
	})
	return //
}

func (row *Row) Len() int64 {
	return int64(len(row.Row.Values))
}

func (cell *Cell) IsBlank() bool {
	if cell == nil || cell.Cell == nil {
		return true
	}
	return cell.GetString() == ""
}

func (cell *Cell) GetString() string {
	if cell == nil || cell.Cell == nil {
		return ""
	}
	return GetCellDataString(cell.Cell)
}

func (cell *Cell) GetAffirm() bool {
	if cell == nil || cell.Cell == nil {
		return false
	}
	val := strings.ToLower(GetCellDataString(cell.Cell))
	return val == "yes" || val == "true" || val == "1"
}
