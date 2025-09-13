// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package sstable_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/clucia/ssdb"
	"github.com/clucia/ssdb/sstable"
	"github.com/stretchr/testify/assert"
)

var ssdbHandle *ssdb.SSDB
var ssTableHandle *sstable.SSTable

func init() {
	spreadsheetID := "<YOUR_SPREADSHEET_ID>"
	serviceAccountKeyFile := "<PATH_TO_YOUR_SERVICE_ACCOUNT_KEY_FILE>"
	credentials, err := os.ReadFile(serviceAccountKeyFile)
	if err != nil {
		log.Fatalf("Unable to read service account key file: %v", err)
		return //
	}
	ctx := context.Background()
	ssdbHandle, err = ssdb.Open(ctx, spreadsheetID, credentials)
	if err != nil {
		log.Fatalf("Open Failed")
	}
	err = ssdbHandle.Loader(ctx)
	if err != nil {
		log.Fatalf("Get Failed")
	}
	ssTableHandle, err = sstable.Open(ssdbHandle, "Config")
	if err != nil {
		log.Fatalf("sstable open Failed")
	}

}

func TestTableLookup(t *testing.T) {
	cell, err := ssTableHandle.Lookup("ClubName", "Global")
	assert.NoError(t, err)
	fmt.Println(cell.GetString())
}
