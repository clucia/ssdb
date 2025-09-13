// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssdb

func (row *Row) SearchH(f func(cell *Cell) bool) *Cell {
	if row == nil {
		return nil
	}
	for i, v := range row.Row.Values {
		cell := &Cell{
			N:     int64(i),
			DB:    row.DB,
			Sheet: row.Sheet,
			Row:   row,
			Cell:  v,
		}
		if f(cell) {
			return cell
		}
	}
	return nil
}

func (sheet *Sheet) SearchV(header bool, f func(row *Row) bool) *Row {
	if sheet == nil {
		return nil
	}
	for i, v := range sheet.Sheet.Data[0].RowData {
		if header && i == 0 {
			continue
		}
		row := &Row{
			N:     int64(i),
			DB:    sheet.DB,
			Sheet: sheet,
			Row:   v,
		}
		if f(row) {
			return row
		}
	}
	return nil
}
