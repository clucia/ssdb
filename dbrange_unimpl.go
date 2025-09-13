// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssdb

type RangeTestCase struct {
	Input       string
	Description string
	// Zero-based coordinates
	UpperLeftRow  int
	UpperLeftCol  int
	LowerRightRow int
	LowerRightCol int
}

var rangeTestCases = []RangeTestCase{
	// Single cell references
	{
		Input:         "A1",
		Description:   "Single cell A1",
		UpperLeftRow:  0,
		UpperLeftCol:  0,
		LowerRightRow: 0,
		LowerRightCol: 0,
	},
	{
		Input:         "Z26",
		Description:   "Single cell Z26",
		UpperLeftRow:  25,
		UpperLeftCol:  25,
		LowerRightRow: 25,
		LowerRightCol: 25,
	},
	{
		Input:         "AA1",
		Description:   "Single cell AA1 (column 26)",
		UpperLeftRow:  0,
		UpperLeftCol:  26,
		LowerRightRow: 0,
		LowerRightCol: 26,
	},
	{
		Input:         "AB100",
		Description:   "Single cell AB100",
		UpperLeftRow:  99,
		UpperLeftCol:  27,
		LowerRightRow: 99,
		LowerRightCol: 27,
	},
	{
		Input:         "AAA1",
		Description:   "Single cell AAA1 (column 702)",
		UpperLeftRow:  0,
		UpperLeftCol:  702,
		LowerRightRow: 0,
		LowerRightCol: 702,
	},

	// Range references
	{
		Input:         "A1:B2",
		Description:   "Simple range A1:B2",
		UpperLeftRow:  0,
		UpperLeftCol:  0,
		LowerRightRow: 1,
		LowerRightCol: 1,
	},
	{
		Input:         "B6:J10",
		Description:   "Range B6:J10",
		UpperLeftRow:  5,
		UpperLeftCol:  1,
		LowerRightRow: 9,
		LowerRightCol: 9,
	},
	{
		Input:         "AA1:ZZ100",
		Description:   "Large range with multi-letter columns",
		UpperLeftRow:  0,
		UpperLeftCol:  26,
		LowerRightRow: 99,
		LowerRightCol: 701,
	},

	// Reversed ranges (should still work)
	{
		Input:         "B2:A1",
		Description:   "Reversed range B2:A1",
		UpperLeftRow:  0,
		UpperLeftCol:  0,
		LowerRightRow: 1,
		LowerRightCol: 1,
	},
	{
		Input:         "J10:B6",
		Description:   "Reversed range J10:B6",
		UpperLeftRow:  5,
		UpperLeftCol:  1,
		LowerRightRow: 9,
		LowerRightCol: 9,
	},

	// Full row references
	{
		Input:         "1:1",
		Description:   "Full row 1",
		UpperLeftRow:  0,
		UpperLeftCol:  0,
		LowerRightRow: 0,
		LowerRightCol: -2, // Special value for "all columns"
	},
	{
		Input:         "5:10",
		Description:   "Full rows 5 to 10",
		UpperLeftRow:  4,
		UpperLeftCol:  0,
		LowerRightRow: 9,
		LowerRightCol: -2,
	},

	// Full column references
	{
		Input:         "A:A",
		Description:   "Full column A",
		UpperLeftRow:  0,
		UpperLeftCol:  0,
		LowerRightRow: -2, // Special value for "all rows"
		LowerRightCol: 0,
	},
	{
		Input:         "B:E",
		Description:   "Full columns B to E",
		UpperLeftRow:  0,
		UpperLeftCol:  1,
		LowerRightRow: -2,
		LowerRightCol: 4,
	},
	{
		Input:         "AA:ZZ",
		Description:   "Full columns AA to ZZ",
		UpperLeftRow:  0,
		UpperLeftCol:  26,
		LowerRightRow: -2,
		LowerRightCol: 701,
	},

	// With sheet names
	{
		Input:         "Sheet1!A1",
		Description:   "Single cell with sheet name",
		UpperLeftRow:  0,
		UpperLeftCol:  0,
		LowerRightRow: 0,
		LowerRightCol: 0,
	},
	{
		Input:         "Sheet1!B6:J10",
		Description:   "Range with sheet name",
		UpperLeftRow:  5,
		UpperLeftCol:  1,
		LowerRightRow: 9,
		LowerRightCol: 9,
	},
	{
		Input:         "MySheet!A1:Z26",
		Description:   "Range with custom sheet name",
		UpperLeftRow:  0,
		UpperLeftCol:  0,
		LowerRightRow: 25,
		LowerRightCol: 25,
	},

	// Sheet names with spaces (quoted)
	{
		Input:         "'Sheet 1'!A1",
		Description:   "Single cell with quoted sheet name",
		UpperLeftRow:  0,
		UpperLeftCol:  0,
		LowerRightRow: 0,
		LowerRightCol: 0,
	},
	{
		Input:         "'My Sheet Name'!B6:J10",
		Description:   "Range with quoted sheet name containing spaces",
		UpperLeftRow:  5,
		UpperLeftCol:  1,
		LowerRightRow: 9,
		LowerRightCol: 9,
	},
	{
		Input:         "'Sheet ''with'' quotes'!A1:B2",
		Description:   "Range with quoted sheet name containing escaped quotes",
		UpperLeftRow:  0,
		UpperLeftCol:  0,
		LowerRightRow: 1,
		LowerRightCol: 1,
	},

	// Edge cases with large numbers
	{
		Input:         "A1000000",
		Description:   "Single cell with large row number",
		UpperLeftRow:  999999,
		UpperLeftCol:  0,
		LowerRightRow: 999999,
		LowerRightCol: 0,
	},

	// Case insensitive
	{
		Input:         "a1:b2",
		Description:   "Lowercase range",
		UpperLeftRow:  0,
		UpperLeftCol:  0,
		LowerRightRow: 1,
		LowerRightCol: 1,
	},
	{
		Input:         "sheet1!a1:b2",
		Description:   "Lowercase range with sheet name",
		UpperLeftRow:  0,
		UpperLeftCol:  0,
		LowerRightRow: 1,
		LowerRightCol: 1,
	},

	// Mixed case
	{
		Input:         "Sheet1!aA1:Bb2",
		Description:   "Mixed case range",
		UpperLeftRow:  0,
		UpperLeftCol:  26,
		LowerRightRow: 1,
		LowerRightCol: 27,
	},

	// Invalid cases
	{
		Input:         "",
		Description:   "Empty string",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         "A",
		Description:   "Just column letter",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         "1",
		Description:   "Just row number",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         "A0",
		Description:   "Invalid row 0",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         "A1:",
		Description:   "Incomplete range (missing end)",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         ":A1",
		Description:   "Incomplete range (missing start)",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         "A1:B2:C3",
		Description:   "Multiple colons",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         "Sheet1!!A1",
		Description:   "Double exclamation",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         "Sheet1!",
		Description:   "Sheet name without range",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         "'Unclosed quote!A1",
		Description:   "Unclosed quote in sheet name",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         "A1.5",
		Description:   "Decimal in row number",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         "123A1",
		Description:   "Numbers before column letters",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         "A-1",
		Description:   "Negative row number",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},

	// Special Google Sheets references
	{
		Input:         "A1:A",
		Description:   "Mixed cell and column reference",
		UpperLeftRow:  0,
		UpperLeftCol:  0,
		LowerRightRow: -2,
		LowerRightCol: 0,
	},
	{
		Input:         "1:B2",
		Description:   "Mixed row and cell reference",
		UpperLeftRow:  0,
		UpperLeftCol:  0,
		LowerRightRow: 1,
		LowerRightCol: 1,
	},

	// Whitespace cases
	{
		Input:         " A1:B2 ",
		Description:   "Range with leading/trailing spaces",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         "A1 : B2",
		Description:   "Range with spaces around colon",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
	{
		Input:         "Sheet 1!A1:B2", // Without quotes
		Description:   "Sheet name with space but no quotes",
		UpperLeftRow:  -1,
		UpperLeftCol:  -1,
		LowerRightRow: -1,
		LowerRightCol: -1,
	},
}
