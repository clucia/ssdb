// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssdb

import (
	"context"
	"sync"
	"time"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/sheets/v4"
)

type SSDB struct {
	sync.Mutex
	ctx           context.Context
	dbTime        time.Time // used for caching
	spreadsheet   *sheets.Spreadsheet
	DocsService   *docs.Service
	SheetsService *sheets.Service
	DriveService  *drive.Service
	SpreadsheetID string
	AppendRows    map[string]int64
}

type Sheet struct {
	N     int64
	DB    *SSDB
	Sheet *sheets.Sheet
}

type Row struct {
	N     int64
	DB    *SSDB
	Sheet *Sheet
	Row   *sheets.RowData
}

type Cell struct {
	N     int64
	DB    *SSDB
	Sheet *Sheet
	Row   *Row
	Cell  *sheets.CellData
}
