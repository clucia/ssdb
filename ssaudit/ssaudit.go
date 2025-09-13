// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssaudit

import (
	"github.com/clucia/ssdb"
	"github.com/clucia/ssdb/sslist"
)

func (sslog *SSAudit) AuditEntry(updater *ssdb.Updater, dat ...any) {
	vals := [][]any{
		dat,
	}
	sslist := (*sslist.SSList)(sslog)
	sslist.AppendBlank(updater, vals)
}
