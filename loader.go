// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssdb

import (
	"context"
	"fmt"
)

func (db *SSDB) Loader(ctx context.Context) (err error) {
	spreadsheet, err := db.SheetsService.Spreadsheets.Get(db.SpreadsheetID).IncludeGridData(true).Context(ctx).Do()
	if err != nil {
		err = fmt.Errorf("unable to get spreadsheet: %v", err)
		return //
	}
	db.spreadsheet = spreadsheet
	return //
}

func (db *SSDB) ReloadDBGet(ctx context.Context) (err error) {
	return db.Loader(ctx)
}
