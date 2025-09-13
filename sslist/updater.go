// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package sslist

import "github.com/clucia/ssdb"

func anySize(_data [][]any) (rows, cols int64) {
	rows = int64(len(_data))
	for _, row := range _data {
		if int64(len(row)) > cols {
			cols = int64(len(row))
		}
	}
	return //
}

func (sslist *SSList) AppendBlank(updater *ssdb.Updater, _data [][]any) {
	inrows, incolumns := anySize(_data)
	dbrange := sslist.GetAppendRange(inrows, incolumns)
	updater.Update(dbrange, _data)
}
