// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package sslist

import (
	"github.com/clucia/ssdb"
)

func (sslist *SSList) GetRange(dbrange *ssdb.DBRange) (_data [][]any) {
	return sslist.Sheet.GetRange(dbrange)
}

func (sslist *SSList) GetAppendRange(rows, columns int64) (dbrange *ssdb.DBRange) {
	line := sslist.GetAppendLine()
	dbrange = sslist.DB.NewDBRange(sslist.sheetName, line, 0, rows, columns)
	return //
}

func (sslist *SSList) GetAppendLine() (newRowN int64) {
	switch {
	case sslist.minAppendLine < 0:
		sslist.Sheet.RowIter(func(row *ssdb.Row) {
			if !row.IsBlank() {
				newRowN = row.N
			}
		})
	default:
		newRowN = sslist.minAppendLine
	}
	sslist.minAppendLine = newRowN + 1
	return newRowN
}
