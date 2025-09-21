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
	"strconv"
	"sync"

	"google.golang.org/api/sheets/v4"
)

type updateItem struct {
	dbRange *DBRange
	olddata [][]any
	newdata [][]any
}

type Updater struct {
	sync.Mutex
	ssdbHandle  *SSDB
	updateQueue []*updateItem
}

func (ssdbHandle *SSDB) NewUpdater() *Updater {
	return &Updater{
		ssdbHandle:  ssdbHandle,
		updateQueue: make([]*updateItem, 0),
	}
}

func (upd *Updater) Update(dbrange *DBRange, newvals [][]any) {
	oldvals := dbrange.sheet.CopyVals(dbrange)
	updtItem := &updateItem{
		dbRange: dbrange,
		olddata: oldvals,
		newdata: newvals,
	}
	upd.Lock()
	upd.updateQueue = append(upd.updateQueue, updtItem)
	upd.Unlock()
}

func (upd *Updater) Len() (l int64) {
	upd.ssdbHandle.Lock()
	defer upd.ssdbHandle.Unlock()

	return int64(len(upd.updateQueue))
}

func (upd *Updater) Sync() (n int, err error) {
	upd.ssdbHandle.Lock()
	defer upd.ssdbHandle.Unlock()

	n = len(upd.updateQueue)
	if n == 0 {
		return // Nothing to sync
	}
	for _, elem := range upd.updateQueue {
		if !elem.dbRange.sheet.CompareVals(elem.dbRange, elem.olddata) {
			err = errors.New("data has changed since the update began")
			return //
		}
	}
	// Create and execute batch update request
	batch := &sheets.BatchUpdateSpreadsheetRequest{}
	for _, update := range upd.updateQueue {
		extents := update.dbRange.sheet.GetExtents()
		growRows, growColumns := extents.NeedsGrowth(update.dbRange)
		if growRows > 0 {
			batch.Requests = append(batch.Requests,
				&sheets.Request{
					AppendDimension: &sheets.AppendDimensionRequest{
						SheetId:   update.dbRange.gridRange.SheetId,
						Dimension: "ROWS",
						Length:    growRows,
					},
				},
			)
		}
		if growColumns > 0 {
			batch.Requests = append(batch.Requests,
				&sheets.Request{
					AppendDimension: &sheets.AppendDimensionRequest{
						SheetId:   update.dbRange.gridRange.SheetId,
						Dimension: "COLUMNS",
						Length:    growColumns,
					},
				},
			)
		}
		rowData := BuildRowdataAny(update.newdata)
		batch.Requests = append(batch.Requests,
			&sheets.Request{
				UpdateCells: &sheets.UpdateCellsRequest{
					Range: update.dbRange.gridRange,
					// show, memberName, text string, data any
					Rows:   rowData,
					Fields: "userEnteredValue,userEnteredFormat.backgroundColor",
				},
			},
		)
	}
	_, err = upd.ssdbHandle.SheetsService.Spreadsheets.BatchUpdate(upd.ssdbHandle.SpreadsheetID, batch).Context(upd.ssdbHandle.ctx).Do()
	if err != nil {
		return 0, fmt.Errorf("unable to batch update spreadsheet: %w", err)
	}

	// Create batch read request for updated ranges
	ranges := make([]string, 0, len(upd.updateQueue))
	for _, update := range upd.updateQueue {
		ranges = append(ranges, upd.ssdbHandle.RangeToString(update.dbRange.gridRange))
	}

	req := upd.ssdbHandle.SheetsService.Spreadsheets.Values.BatchGet(upd.ssdbHandle.SpreadsheetID).
		Ranges(ranges...).
		ValueRenderOption("FORMATTED_VALUE")

	// Execute the batch read request
	resp, err := req.Context(upd.ssdbHandle.ctx).Do()
	if err != nil {
		err = fmt.Errorf("unable to retrieve data from sheet: %w", err)
		return //
	}

	// Merge the results
	err = upd.ssdbHandle.Merge(resp)
	if err != nil {
		err = fmt.Errorf("unable to merge data from sheet: %w", err)
		return //
	}
	upd.updateQueue = make([]*updateItem, 0)
	return //
}

func BuildRowdataAny(data [][]any) (rowData []*sheets.RowData) {
	xdim, ydim := getDimsAny(data)
	rowData = []*sheets.RowData{}
	for ypos := 0; ypos < ydim; ypos++ {
		rowData = append(rowData, &sheets.RowData{
			Values: make([]*sheets.CellData, xdim),
		})
		for xpos := 0; xpos < len(data[ypos]); xpos++ {
			var fld string
			switch fld0 := data[ypos][xpos].(type) {
			case string:
				fld = fld0
			default:
				fld = fmt.Sprint(fld0)
			}
			switch {
			case isBlank(fld):
				rowData[ypos].Values[xpos] = &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						BackgroundColor: AttnColor,
					},
				}
			case isNumeric(fld):
				numval, _ := strconv.ParseFloat(fld, 64)
				rowData[ypos].Values[xpos] = &sheets.CellData{
					UserEnteredValue: &sheets.ExtendedValue{
						NumberValue: &numval,
					},
					UserEnteredFormat: &sheets.CellFormat{
						BackgroundColor: AttnColor,
					},
					FormattedValue: fmt.Sprint(numval),
				}
			case isString(fld):
				rowData[ypos].Values[xpos] = &sheets.CellData{
					UserEnteredValue: &sheets.ExtendedValue{
						StringValue: &fld,
					},
					UserEnteredFormat: &sheets.CellFormat{
						BackgroundColor: AttnColor,
					},
					FormattedValue: fld,
				}
			}
		}
	}
	return //
}

func BuildRowdata(data [][]string) (rowData []*sheets.RowData) {
	xdim, ydim := getDims(data)
	rowData = []*sheets.RowData{}
	for ypos := 0; ypos < ydim; ypos++ {
		rowData = append(rowData, &sheets.RowData{
			Values: make([]*sheets.CellData, xdim),
		})
		for xpos := 0; xpos < len(data[ypos]); xpos++ {
			switch {
			case isBlank(data[ypos][xpos]):
				rowData[ypos].Values[xpos] = &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						BackgroundColor: AttnColor,
					},
				}
			case isNumeric(data[ypos][xpos]):
				numval, _ := strconv.ParseFloat(data[ypos][xpos], 64)
				rowData[ypos].Values[xpos] = &sheets.CellData{
					UserEnteredValue: &sheets.ExtendedValue{
						NumberValue: &numval,
					},
					UserEnteredFormat: &sheets.CellFormat{
						BackgroundColor: AttnColor,
					},
					FormattedValue: fmt.Sprint(numval),
				}
			case isString(data[ypos][xpos]):
				rowData[ypos].Values[xpos] = &sheets.CellData{
					UserEnteredValue: &sheets.ExtendedValue{
						StringValue: &data[ypos][xpos],
					},
					UserEnteredFormat: &sheets.CellFormat{
						BackgroundColor: AttnColor,
					},
					FormattedValue: data[ypos][xpos],
				}
			}
		}
	}
	return //
}

func getDimsAny(data [][]any) (x, y int) {
	x = -1
	y = -1
	for j, v := range data {
		for i, _ := range v {
			if i > x {
				x = i
			}
		}
		if j > y {
			y = j
		}
	}
	return x + 1, y + 1
}

func getDims(data [][]string) (x, y int) {
	x = -1
	y = -1
	for j, v := range data {
		for i, _ := range v {
			if i > x {
				x = i
			}
		}
		if j > y {
			y = j
		}
	}
	return x + 1, y + 1
}

func isBlank(s string) bool {
	return len(s) == 0
}

func isString(s string) bool {
	if isNumeric(s) {
		return false
	}
	return len(s) > 0
}

func isNumeric(s string) bool {
	switch {
	case len(s) == 0:
		return false
	case s == "$":
		return false
	case s[0] == '$':
		s = s[1:]
	}
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if (c < '0' || c > '9') && c != '.' {
			return false
		}
	}
	return true
}

var AttnColor = &sheets.Color{
	Red:   1.0,
	Green: 1.0,
	Blue:  0.4,
}

func AnyIfy(in [][]string) (out [][]any) {
	for _, row := range in {
		line := []any{}
		for _, item := range row {
			line = append(line, item)
		}
		out = append(out, line)
	}
	return //
}

func DeAnyIfy(in [][]any) (out [][]string) {
	for _, row := range in {
		line := []string{}
		for _, item := range row {
			switch it := item.(type) {
			case string:
				line = append(line, it)
			default:
				line = append(line, fmt.Sprint(item))
			}
		}
		out = append(out, line)
	}
	return //
}

func (db *SSDB) Merge(rresp *sheets.BatchGetValuesResponse) error {
	if rresp == nil {
		return fmt.Errorf("response is nil")
	}
	for _, vr := range rresp.ValueRanges {
		if err := db.mergeValueRange(vr); err != nil {
			return fmt.Errorf("failed to merge range %s: %w", vr.Range, err)
		}
	}
	return nil
}

func (db *SSDB) mergeValueRange(vr *sheets.ValueRange) error {
	rng, err := db.TextToSheetsRange(vr.Range)
	if err != nil {
		return err
	}

	sheet := db.FindSheet(rng)
	if sheet == nil {
		return fmt.Errorf("sheet not found for range: %s", vr.Range)
	}

	if err := db.ensureSheetCapacity(sheet, rng); err != nil {
		return err
	}

	return db.updateCells(sheet, rng, vr.Values)
}

func (db *SSDB) updateCells(sheet *Sheet, rng *sheets.GridRange, values [][]interface{}) error {
	for rowIdx, rowValues := range values {
		actualRow := int(rng.StartRowIndex) + rowIdx

		if actualRow >= len(sheet.Sheet.Data[0].RowData) {
			continue // Skip if out of bounds
		}

		for colIdx, cellValue := range rowValues {
			actualCol := int(rng.StartColumnIndex) + colIdx

			if actualCol >= len(sheet.Sheet.Data[0].RowData[actualRow].Values) {
				continue // Skip if out of bounds
			}

			// Safe type conversion
			cellStr, ok := cellValue.(string)
			if !ok {
				cellStr = fmt.Sprintf("%v", cellValue)
			}

			existingCell := sheet.Sheet.Data[0].RowData[actualRow].Values[actualCol]
			sheet.Sheet.Data[0].RowData[actualRow].Values[actualCol] = genCell(existingCell, cellStr)
		}
	}
	return nil
}

func (db *SSDB) ensureSheetCapacity(sheet *Sheet, rng *sheets.GridRange) error {
	// Validate inputs
	if sheet == nil {
		return fmt.Errorf("sheet is nil")
	}
	if rng == nil {
		return fmt.Errorf("range is nil")
	}
	// Initialize sheet data structure if needed
	if sheet.Sheet.Data == nil {
		sheet.Sheet.Data = []*sheets.GridData{{}}
	}
	if len(sheet.Sheet.Data) == 0 {
		sheet.Sheet.Data = append(sheet.Sheet.Data, &sheets.GridData{})
	}
	gridData := sheet.Sheet.Data[0]

	// Calculate required dimensions
	requiredRows := int(rng.EndRowIndex)
	requiredCols := int(rng.EndColumnIndex)

	// Ensure we have enough rows
	currentRows := len(gridData.RowData)
	if currentRows < requiredRows {
		// Pre-allocate the slice to avoid multiple reallocations
		newRows := make([]*sheets.RowData, requiredRows)
		copy(newRows, gridData.RowData)

		// Initialize new rows
		for i := currentRows; i < requiredRows; i++ {
			newRows[i] = &sheets.RowData{
				Values: make([]*sheets.CellData, 0),
			}
		}
		gridData.RowData = newRows
	}

	// Ensure each row has enough columns
	for rowIdx := 0; rowIdx < requiredRows; rowIdx++ {
		row := gridData.RowData[rowIdx]
		if row == nil {
			row = &sheets.RowData{Values: make([]*sheets.CellData, 0)}
			gridData.RowData[rowIdx] = row
		}

		currentCols := len(row.Values)
		if currentCols < requiredCols {
			// Pre-allocate the slice
			newCells := make([]*sheets.CellData, requiredCols)
			copy(newCells, row.Values)

			// Initialize new cells
			for colIdx := currentCols; colIdx < requiredCols; colIdx++ {
				newCells[colIdx] = &sheets.CellData{}
			}
			row.Values = newCells
		}
	}

	return nil
}

func (db *SSDB) FindSheet(rng *sheets.GridRange) *Sheet {
	for _, sheet := range db.spreadsheet.Sheets {
		if sheet.Properties.SheetId == rng.SheetId {
			return &Sheet{
				DB:    db,
				Sheet: sheet,
			}
		}
	}
	return nil
}

func genCell(cell *sheets.CellData, cal string) (res *sheets.CellData) {

	if cell.UserEnteredValue == nil {
		cell.UserEnteredValue = &sheets.ExtendedValue{}
	}
	switch {
	case isBlank(cal):
	case isNumeric(cal):
		numval, _ := strconv.ParseFloat(cal, 64)
		cell.UserEnteredValue.NumberValue = &numval
		cell.FormattedValue = fmt.Sprint(numval)
	case isString(cal):
		cell.UserEnteredValue.StringValue = &cal
		cell.FormattedValue = cal
	}
	if cell.UserEnteredFormat == nil {
		cell.UserEnteredFormat = &sheets.CellFormat{}
	}
	cell.UserEnteredFormat.BackgroundColor = AttnColor
	return cell
}
