// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssdb_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/clucia/ssdb"
)

var ssdbHandle *ssdb.SSDB

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
}

func TestSheetIter(t *testing.T) {
	ssdbHandle.SheetIter(func(name string, sheet *ssdb.Sheet) {
		fmt.Print("\nsheet ", name, ": ============================== \n")
		sheet.RowIter(func(row *ssdb.Row) {
			fmt.Print("\nRow", row.N, ": ")
			row.CellIter(func(cell *ssdb.Cell) {
				fmt.Print("cell ", cell.N, ", ")
			})
		})
	})
}
