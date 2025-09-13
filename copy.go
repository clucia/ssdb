// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssdb

func (sheet *Sheet) CompareVals(dbrange *DBRange, vals [][]any) (res bool) {
	row0, col0, row1, col1 := dbrange.gridRange.StartRowIndex, dbrange.gridRange.StartColumnIndex,
		dbrange.gridRange.EndRowIndex, dbrange.gridRange.EndColumnIndex
	_, _, _, _ = row0, col0, row1, col1

	res = true
	sheet.RowIter(func(row *Row) {
		if !res || row.N < row0 || row.N > row1 {
			return
		}
		row.CellIter(func(cell *Cell) {
			if !res || cell.N < col0 || cell.N > col1 {
				return
			}
			if vals[int(row.N-row0)][int(cell.N-col0)].(string) != cell.GetString() {
				res = false
				return
			}
		})
	})
	return //
}

func (sheet *Sheet) CopyVals(dbrange *DBRange) (vals [][]any) {
	row0, col0, row1, col1 := dbrange.gridRange.StartRowIndex, dbrange.gridRange.StartColumnIndex,
		dbrange.gridRange.EndRowIndex, dbrange.gridRange.EndColumnIndex
	_, _, _, _ = row0, col0, row1, col1

	sheet.RowIter(func(row *Row) {
		if row.N < row0 || row.N > row1 {
			return
		}
		row.CellIter(func(cell *Cell) {
			if cell.N < col0 || cell.N > col1 {
				return
			}
			l := len(vals)
			add := int(row.N-row0) - l + 1
			for i := 0; i < add; i++ {
				vals = append(vals, []any{})
			}
			vals[row.N-row0] = append(vals[row.N-row0], cell.GetString())
		})
	})
	return //
}
