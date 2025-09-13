// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package sschat

import (
	"github.com/clucia/ssdb"
	"github.com/clucia/ssdb/sslist"
)

func (sslog *SSChat) AppendChat(updater *ssdb.Updater, dat ...any) {
	vals := [][]any{
		dat,
	}
	sslist := (*sslist.SSList)(sslog)
	sslist.AppendBlank(updater, vals)
}

func (sslog *SSChat) GetChat() (res []map[string]string) {
	sslist := (*sslist.SSList)(sslog)

	var hdrRow *ssdb.Row
	sslist.Sheet.RowIter(func(row *ssdb.Row) {
		if row.N == 0 {
			hdrRow = row
			return
		}
		ent := map[string]string{}
		row.CellIter(func(cell *ssdb.Cell) {
			hdr := hdrRow.GetCellN(cell.N).GetString()
			ent[hdr] = cell.GetString()
		})
		res = append(res, ent)
	})
	return //
}
