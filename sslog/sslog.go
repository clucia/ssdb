// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package sslog

import (
	"context"
	"time"

	"github.com/clucia/ssdb"
	"github.com/clucia/ssdb/sslist"
)

func (sslog *SSLog) Log(updater *ssdb.Updater, dat ...any) {
	line := []any{time.Now().Format(time.RFC3339)}
	line = append(line, dat...)
	vals := [][]any{line}
	sslist := (*sslist.SSList)(sslog)
	sslist.AppendBlank(updater, vals)
	updater.Sync()
}

func (sslog *SSLog) LogErr(ctx context.Context, dat ...any) {
	sslist := (*sslist.SSList)(sslog)
	updater := sslist.DB.NewUpdater()
	line := []any{}
	line = append(line, time.Now().Format(time.RFC3339))
	line = append(line, "ERROR")
	line = append(line, dat...)
	vals := [][]any{line}
	sslist.AppendBlank(updater, vals)
	updater.Sync()
}

func (sslog *SSLog) LogWithData(ctx context.Context, dat ...any) {
	sslist := (*sslist.SSList)(sslog)
	updater := sslist.DB.NewUpdater()
	line := []any{}
	line = append(line, time.Now().Format(time.RFC3339))
	line = append(line, "DATA")
	line = append(line, dat...)
	vals := [][]any{line}
	sslist.AppendBlank(updater, vals)
	updater.Sync()
}
