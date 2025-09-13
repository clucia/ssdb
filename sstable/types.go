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

type SSTable struct {
	DB        *ssdb.SSDB
	sheetName string
	sheet     *ssdb.Sheet
}

var ErrSheetNotFound = errors.New("sheet not found")

func Open(db *ssdb.SSDB, match string) (sstable *SSTable, err error) {
	sstable = &SSTable{
		DB: db,
	}
	var found bool
	db.SheetIter(func(sheetname string, sheet *ssdb.Sheet) {
		if found {
			return
		}
		if sheetname == match {
			found = true
			sstable.sheet = sheet
			sstable.sheetName = sheetname
		}
	})
	if !found {
		sstable, err = nil, ErrSheetNotFound
	}
	return //
}
