// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssdb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdaterGrow(t *testing.T) {
	updater := ssdbHandle.NewUpdater()
	dbRange := ssdbHandle.NewDBRangeFromSymbolicRange("Config!J5:K6")
	newvals := [][]any{
		{"J5", "K5"},
		{"J6", "K6"},
	}
	updater.Update(dbRange, newvals)

	dbRange = ssdbHandle.NewDBRangeFromSymbolicRange("Config!A55:B56")
	newvals = [][]any{
		{"A55", "B55"},
		{"A56", "B56"},
	}
	updater.Update(dbRange, newvals)

	dbRange = ssdbHandle.NewDBRangeFromSymbolicRange("Config!L75:M76")
	newvals = [][]any{
		{"L75", "M75"},
		{"L76", "M76"},
	}
	updater.Update(dbRange, newvals)

	err := updater.Sync()
	assert.NoError(t, err, "extending range failed")
}

func TestUpdater(t *testing.T) {
	updater := ssdbHandle.NewUpdater()
	dbRange := ssdbHandle.NewDBRangeFromSymbolicRange("Config!C17:D29")
	newvals := [][]any{
		{"foo", "bar"},
		{"foo", "bar"},
		{"foo", "bar"},
		{"foo", "bar"},
		{"foo", "bar"},
		{"foo", "bar"},
		{"bax", "lbj"},
		{"foo", "bar"},
		{"foo", "bar"},
		{"foo", "bar"},
		{"foo", "bar"},
		{"foo", "bar"},
		{"foo", "bar"},
	}
	updater.Update(dbRange, newvals)
	dbRange = ssdbHandle.NewDBRangeFromSymbolicRange("Config!C12")
	newvals = [][]any{
		{"TRUE"},
	}
	updater.Update(dbRange, newvals)
	updater.Sync()

	_ = updater
}
