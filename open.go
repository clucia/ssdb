// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssdb

import (
	"context"
	"log"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func Open(
	ctx context.Context,
	spreadsheetID string,
	credentials []byte,
) (db *SSDB, err error) {
	configAPI, err := google.JWTConfigFromJSON(
		credentials,
		sheets.DriveScope,
		sheets.DriveFileScope,
		sheets.SpreadsheetsScope,
	)
	if err != nil {
		log.Fatalf("Unable to parse service account key file to config: %v", err)
		return //
	}
	client := configAPI.Client(ctx)
	// Create a new Doc docService.
	docsService, err := docs.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
		return //
	}
	// Create a new Sheets sheetsService.
	sheetsService, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
		return //
	}
	// Create the Drive service
	driveService, err := drive.NewService(ctx, option.WithCredentials(&google.Credentials{
		JSON: credentials,
	}))
	if err != nil {
		log.Fatalf("Unable to create drive service: %v", err)
		return //
	}
	log.Printf("GetLive")
	db = &SSDB{
		ctx:           ctx,
		dbTime:        time.Now(),
		SheetsService: sheetsService,
		DriveService:  driveService,
		DocsService:   docsService,
		SpreadsheetID: spreadsheetID,
	}

	return //
}
