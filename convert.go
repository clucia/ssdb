// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssdb

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/api/sheets/v4"
)

func GenRange(sheetID, minx, maxx, miny, maxy int64) (rng *sheets.GridRange) {
	return &sheets.GridRange{
		SheetId:          sheetID,
		StartColumnIndex: minx,
		EndColumnIndex:   maxx,
		StartRowIndex:    miny,
		EndRowIndex:      maxy,
	}
}

func GenRangeString(page string, minx, maxx, miny, maxy int) string {
	maxx--
	maxy--
	f := func(n int) (res string) {
		m, o := byte(n%26), byte(n/26)
		if o > 0 {
			return string([]byte{'A' + o, 'A' + m})
		}
		return string([]byte{'A' + m})
	}
	return fmt.Sprintf("%s!%s%d:%s%d", page, f(minx), miny, f(maxx), maxy)

}

// TextToSheetsRange converts a text representation back to a *sheets.Range
// Handles formats like:
// - Sheet1!A1:B2
// - 'Sheet Name with spaces'!A1:C3
// - A1:D5 (assumes current sheet)
func (db *SSDB) TextToSheetsRange(text string) (result *sheets.GridRange, err error) {
	result = &sheets.GridRange{
		SheetId: -1,
	}
	spl := strings.Split(text, "!")
	if len(spl) < 2 {
		return nil, errors.New("bad first split")
	}
	for _, sheet := range db.spreadsheet.Sheets {
		if sheet.Properties.Title == spl[0] {
			result.SheetId = sheet.Properties.SheetId
		}
	}
	if result.SheetId < 0 {
		return nil, errors.New("bad sheet id")
	}

	rangeRegex := regexp.MustCompile(`^([A-Z]+)(\d+)(?::([A-Z]+)(\d+))?$`)
	rangeMatches := rangeRegex.FindStringSubmatch(spl[1])

	if len(rangeMatches) == 0 {
		return nil, fmt.Errorf("invalid range format: %s", spl[1])
	}

	// Start column and row
	result.StartColumnIndex = columnLetterToIndex(rangeMatches[1])
	var startRow int
	fmt.Sscanf(rangeMatches[2], "%d", &startRow)
	result.StartRowIndex = int64(startRow - 1) // Convert to 0-based index

	// If there's an end range
	if rangeMatches[3] != "" {
		result.EndColumnIndex = columnLetterToIndex(rangeMatches[3]) + 1 // +1 because end is exclusive
		var endRow int
		fmt.Sscanf(rangeMatches[4], "%d", &endRow)
		result.EndRowIndex = int64(endRow) // Convert to 0-based and +1 for exclusive
	} else {
		// Single cell, make range 1x1
		result.EndColumnIndex = result.StartColumnIndex + 1
		result.EndRowIndex = result.StartRowIndex + 1
	}

	return result, nil
}

// columnLetterToIndex converts a column letter like "A", "B", "AA", etc. to a 0-based index
func columnLetterToIndex(letters string) int64 {
	result := int64(0)
	for _, letter := range letters {
		result = result*26 + int64(letter-'A'+1)
	}
	return result - 1 // Convert to 0-based index
}

func (db *SSDB) RangeToString(rng *sheets.GridRange) (rngst string) {
	if rng == nil {
		return ""
	}
	page := ""
	for _, sheet := range db.spreadsheet.Sheets {
		if rng.SheetId == sheet.Properties.SheetId {
			page = sheet.Properties.Title
		}
	}
	// Convert start and end columns to letters
	startCol := columnIndexToLetter(int(rng.StartColumnIndex))
	endCol := columnIndexToLetter(int(rng.EndColumnIndex - 1)) // Subtract 1 because GridRange is end-exclusive

	// Adjust row indices (add 1 because sheets rows are 1-based in A1 notation)
	startRow := rng.StartRowIndex + 1
	endRow := rng.EndRowIndex // No need to subtract 1 as we'll use gr.EndRowIndex directly

	// Format as A1 notation
	return fmt.Sprintf("%s!%s%d:%s%d", page, startCol, startRow, endCol, endRow)
}

func RangeToString(page string, rng *sheets.GridRange) (rngst string) {
	f := func(n int) (res string) {
		m, o := byte(n%26), byte(n/26)
		if o > 0 {
			return string([]byte{'A' + o, 'A' + m})
		}
		return string([]byte{'A' + m})
	}
	return fmt.Sprintf("%s!%s%d:%s%d",
		page,

		f(int(rng.StartColumnIndex)), rng.StartRowIndex,
		f(int(rng.EndColumnIndex)), rng.EndRowIndex,
	)
}

func RangeFromString(sheetID int64, rngst string) *sheets.GridRange {
	bang := strings.Split(rngst, "!")

	page := ""
	if len(bang) == 2 {
		page = bang[0]
	}
	fmt.Println(page)
	corn := strings.Split(bang[1], ":")
	xmin, xmax := int64(0), int64(0)
	f := func(in string) (x, y int64) {
		x, y = int64(0), int64(0)
		for _, c := range in {
			switch {
			case c >= '0' && c <= '9':
				y = y*10 + int64(c-'0')
			case c >= 'A' && c <= 'Z':
				x = x*26 + int64(c-'A')
			}
		}
		return //
	}
	xmin, xmax, ymin, ymax := int64(0), int64(0), int64(0), int64(0)
	if len(corn) == 2 {
		xmin, ymin = f(corn[0])
		xmax, ymax = f(corn[1])
		xmax++
		ymax++
	} else {
		xmin, ymin = f(corn[0])
		xmax, ymax = xmin+1, ymin+1
	}
	return &sheets.GridRange{
		SheetId:          sheetID,
		StartRowIndex:    ymin,
		EndRowIndex:      ymax,
		StartColumnIndex: xmin,
		EndColumnIndex:   xmax,
	}
}

func GetCellDataString(cell *sheets.CellData) (val string) {
	val = cell.FormattedValue
	if cell.UserEnteredValue != nil &&
		cell.UserEnteredValue.StringValue != nil &&
		*cell.UserEnteredValue.StringValue != "" {
		val = *cell.UserEnteredValue.StringValue
	}
	return //
}

// Utility function to convert column index to letter (0=A, 1=B, etc.)
func columnIndexToLetter(index int) string {
	result := ""
	for index >= 0 {
		result = string(rune('A'+index%26)) + result
		index = index/26 - 1
	}
	return result
}
