// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package sstable

import (
	"fmt"

	"github.com/clucia/ssdb"
)

// ctx context.Context, gsend, tsend, emailAddr, phoneNumber, textTempl, emailTempl, templData map[string]any) {
func (sstable *SSTable) Iter(f func(_map map[string]any)) {
	var hdrrow *ssdb.Row
	var enableHdrCell *ssdb.Cell
	var _map map[string]any

	sstable.sheet.RowIter(func(row *ssdb.Row) {
		_map = nil
		switch {
		case row.N == 0:
			enableHdrCell, _ = row.GetCellByName("Enable") // No enable header implies everything is enabled // TODO
			hdrrow = row
			return
		case row.IsBlank():
			return
		// if we have an enable colum and enable is false on this row, skip it
		case enableHdrCell != nil && !row.GetCellN(enableHdrCell.N).GetAffirm():
			// skip
		default:
			row.CellIter(func(cell *ssdb.Cell) {
				if _map == nil {
					_map = make(map[string]any)
				}
				switch {
				case cell.N < hdrrow.Len():
					hdr := hdrrow.GetCellN(cell.N).GetString()
					val := row.GetCellN(cell.N).GetString()
					_map[hdr] = val
				default:
					hdr := fmt.Sprintf("Unnamed(%d)", cell.N)
					val := row.GetCellN(cell.N).GetString()
					_map[hdr] = val
				}
			})
			f(_map)
		}
	})
}

func (sstable *SSTable) Get2DMap(keyColumn string) (_map map[string]map[string]any) {
	var hdrrow *ssdb.Row
	var keyFlag *ssdb.Cell

	_map = nil
	sstable.sheet.RowIter(func(row *ssdb.Row) {
		switch {
		case row.N == 0:
			hdrrow = row
			keyFlag, _ = row.GetCellByName(keyColumn) // TODO
			return
		case row.IsBlank():
			return
		// if we have an enable colum and enable is false on this row, skip it
		default:
			var key, k, v string
			if keyFlag == nil {
				key = fmt.Sprintf("UnamedRow%d", row.N)
			} else {
				keyCell, _ := row.GetCellByName(keyColumn) // TODO)
				key = keyCell.GetString()
			}
			row.CellIter(func(cell *ssdb.Cell) {
				switch {
				case cell.N < hdrrow.Len():
					k = hdrrow.GetCellN(cell.N).GetString()
				default:
					k = fmt.Sprintf("UnnamedColumn(%d)", cell.N)
				}
				v = row.GetCellN(cell.N).GetString()
				if _map == nil {
					_map = make(map[string]map[string]any)
				}
				if _map[key] == nil {
					_map[key] = make(map[string]any)
				}
				_map[key][k] = v
			})
		}
	})
	return //
}

func (sstable *SSTable) GetHeaders() (headers []string) {
	var hdrrow *ssdb.Row

	hdrrow = sstable.sheet.GetRowN(0)
	hdrrow.CellIter(func(cell *ssdb.Cell) {
		headers = append(headers, cell.GetString())
	})
	return //
}

func (sstable *SSTable) ListColumn(N int64) (list []string) {
	sstable.sheet.RowIter(func(row *ssdb.Row) {
		switch {
		case row.N == 0:
			return
		case row.Len() <= N:
			list = append(list, "Short Row")
		default:
			list = append(list, row.GetCellN(N).GetString())
		}
	})
	return //
}

func (sstable *SSTable) ListColumnByName(name string) (list []string) {
	var hdrrow *ssdb.Row
	hdrrow = sstable.sheet.GetRowN(0)
	colN := int64(-1)
	var err bool
	hdrrow.CellIter(func(cell *ssdb.Cell) {
		switch {
		case err:
			return
		case colN != -1 && cell.GetString() == name:
			err = true // shouldn't match more than one
		case cell.GetString() == name:
			colN = cell.N
		default:
		}
	})
	return sstable.ListColumn(colN)
}
