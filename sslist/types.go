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

func Open(db *ssdb.SSDB, match string) (sslst *SSList, err error) {
	sslst = &SSList{
		DB: db,
	}
	var found bool
	db.SheetIter(func(sheetname string, sheet *ssdb.Sheet) {
		if found {
			return
		}
		if sheetname == match {
			found = true
			sslst.Sheet = sheet
			sslst.sheetName = sheetname
		}
	})
	if !found {
		sslst, err = nil, ErrSheetNotFound
	}
	return //
}
