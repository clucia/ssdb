// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package sslog

import (
	"errors"

	"github.com/clucia/ssdb"
	"github.com/clucia/ssdb/sslist"
)

type SSLog sslist.SSList

var ErrSheetNotFound = errors.New("sheet not found")

func Open(ssdb *ssdb.SSDB, matchSheet string) (sslog *SSLog, err error) {
	sslist, err := sslist.Open(ssdb, matchSheet)
	return (*SSLog)(sslist), err
}
