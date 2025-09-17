// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssaudit

import (
	"time"

	"github.com/clucia/ssdb"
	"github.com/clucia/ssdb/sslist"
)

func (sslog *SSAudit) AuditUpdate(updater *ssdb.Updater, dat ...any) {
	line := []any{
		time.Now().Format(time.RFC3339),
	}
	line = append(line, dat...)
	vals := [][]any{line}
	sslist := (*sslist.SSList)(sslog)
	sslist.AppendBlank(updater, vals)
}

func (sslog *SSAudit) AuditEntry(updater *ssdb.Updater, dat ...any) {
	sslog.AuditUpdate(updater, dat...)
	updater.Sync()
}
