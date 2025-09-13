// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package sslog

import (
	"context"

	"github.com/clucia/ssdb"
	"github.com/clucia/ssdb/sslist"
)

func (sslog *SSLog) Log(updater *ssdb.Updater, dat ...any) {
	vals := [][]any{
		dat,
	}
	sslist := (*sslist.SSList)(sslog)
	sslist.AppendBlank(updater, vals)
}

func (sslog *SSLog) LogErr(ctx context.Context, dat ...any) {
	sslist := (*sslist.SSList)(sslog)
	updater := sslist.DB.NewUpdater()
	vals := [][]any{
		append([]any{"Error"}, dat...),
	}
	sslist.AppendBlank(updater, vals)
	updater.Sync()
}

func (sslog *SSLog) LogWithData(ctx context.Context, dat ...any) {
	sslist := (*sslist.SSList)(sslog)
	updater := sslist.DB.NewUpdater()
	vals := [][]any{
		append([]any{"Data"}, dat...),
	}
	sslist.AppendBlank(updater, vals)
	updater.Sync()
}
