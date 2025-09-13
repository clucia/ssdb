// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssdb

func (db *SSDB) SheetIter(f func(sheetname string, sheet *Sheet)) {
	for i, Gsheet := range db.spreadsheet.Sheets {
		sheetName := Gsheet.Properties.Title
		sheet := &Sheet{
			N:     int64(i),
			DB:    db,
			Sheet: Gsheet,
		}
		f(sheetName, sheet)
	}
}

func (sheet *Sheet) RowIter(f func(row *Row)) {
	for i, row := range sheet.Sheet.Data[0].RowData {
		f(&Row{
			N:     int64(i),
			DB:    sheet.DB,
			Sheet: sheet,
			Row:   row,
		})
	}
}

func (row *Row) CellIter(f func(cell *Cell)) {
	for N, cell := range row.Row.Values {
		f(&Cell{
			N:     int64(N),
			DB:    row.DB,
			Sheet: row.Sheet,
			Row:   row,
			Cell:  cell,
		})
	}
}

func (sheet *Sheet) GetID() int64 {
	return sheet.Sheet.Properties.SheetId
}
