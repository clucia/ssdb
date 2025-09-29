// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package sslist

import (
	"errors"

	"github.com/clucia/ssdb"
)

type SSList struct {
	DB            *ssdb.SSDB
	sheetName     string
	Sheet         *ssdb.Sheet
	minAppendLine int64
}

var ErrSheetNotFound = errors.New("sheet not found")
var ErrCantAppend = errors.New("cannot append")

func Open(db *ssdb.SSDB, match string) (sslst *SSList, err error) {
	sslst = &SSList{
		DB: db,
	}
	sheet := db.SheetLookup(match)
	if sheet == nil {
		sslst, err = nil, ErrSheetNotFound
	} else {
		sslst = &SSList{
			DB:        db,
			sheetName: match,
			Sheet:     sheet,
		}
	}
	minAppend := int64(-1)
	sheet.RowIter(func(row *ssdb.Row) {
		if !row.IsBlank() {
			minAppend = row.N
		}
	})
	if minAppend < 0 {
		sslst, err = nil, ErrCantAppend
		return
	}
	sslst.minAppendLine = minAppend + 1
	return //
}
