// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package sstable

import (
	"errors"

	"github.com/clucia/ssdb"
)

var ErrLookupFailed = errors.New("lookup failed")
var ErrDuplicateRowKey = errors.New("duplicate row key")
var ErrDuplicateColumnKey = errors.New("duplicate column key")

func (sstbl *SSTable) HSearch(colMatch string) (foundCell *ssdb.Cell) {
	var hdrrow *ssdb.Row
	sstbl.sheet.RowIter(func(row *ssdb.Row) {
		if row.N == 0 {
			hdrrow = row
		}
	})
	var err bool
	hdrrow.CellIter(func(cell *ssdb.Cell) {
		switch {
		case err:
			return
		case foundCell != nil && cell.GetString() == colMatch:
			err = true
		case cell.GetString() == colMatch:
			foundCell = cell
		}
	})
	if err {
		return nil
	}
	return foundCell
}

func (sstbl *SSTable) VSearch(colMatch string, colValue string) (foundrow *ssdb.Row) {
	var err bool
	cell := sstbl.HSearch(colMatch)
	sstbl.sheet.RowIter(func(row *ssdb.Row) {
		if row.N == 0 {
			return
		}
		testcell := row.GetCellN(cell.N)
		switch {
		case foundrow != nil && testcell.GetString() == colValue:
			err = true
		case testcell.GetString() == colValue:
			foundrow = row
		}
	})
	if err {
		return nil
	}
	return //
}

func (sstbl *SSTable) Lookup(rowMatch, colMatch string) (cell *ssdb.Cell, err error) {
	var foundRow *ssdb.Row
	sstbl.sheet.RowIter(func(row *ssdb.Row) {
		if err != nil {
			return
		}
		if row.GetCellN(0).GetString() == rowMatch {
			if foundRow != nil {
				err = ErrDuplicateRowKey
			}
			foundRow = row
		}
	})
	if foundRow == nil {
		err = ErrLookupFailed
		return
	}
	var foundCell *ssdb.Cell
	hdrrow := sstbl.sheet.GetRowN(0)
	hdrrow.CellIter(func(cell *ssdb.Cell) {
		if cell.GetString() == colMatch {
			if foundCell != nil {
				err = ErrDuplicateColumnKey
				return
			}
			foundCell = cell
		}
	})
	if foundCell == nil {
		err = ErrLookupFailed
		return //
	}
	cell = foundRow.GetCellN(foundCell.N)
	return //
}

func (sstbl *SSTable) GetKeys() (rowKeys, colKeys []string) {
	sstbl.sheet.RowIter(func(row *ssdb.Row) {
		if row.N == 0 {
			return
		}
		if row.Len() == 0 {
			return
		}
		cell := row.GetCellN(0)
		if cell == nil {
			return
		}
		cellN := cell.GetString()
		if cellN == "" {
			return
		}
		rowKeys = append(rowKeys, row.GetCellN(0).GetString())
	})
	hdrrow := sstbl.sheet.GetRowN(0)
	hdrrow.CellIter(func(cell *ssdb.Cell) {
		colKeys = append(colKeys, cell.GetString())
	})
	return //
}
