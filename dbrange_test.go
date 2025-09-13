// Copyright (c) 2025 Way To Go LLC. All rights reserved.
//
// This file is part of SSDB (Spreadsheet Database).
//
// Licensed under the MIT License. See LICENSE file in the project root
// for full license information.
package ssdb_test

import (
	"fmt"
	"testing"

	"github.com/clucia/ssdb"
	"github.com/stretchr/testify/assert"
)

func TestNewDBRange(t *testing.T) {
	var row, col, rows, cols int64

	sheetname := "log"
	row = 27
	col = 0
	rows = 1
	cols = 4

	dbrange := ssdbHandle.NewDBRange(sheetname, row, col, rows, cols)
	fmt.Println(dbrange)

}

func TestParseSymbolicRange(t *testing.T) {
	tbl := []struct {
		in     string
		corner bool // false is NW, true is SE
		row    int
		col    int
		name   string
	}{
		{in: "b7", corner: false, row: 6, col: 1, name: "testname"},
		{in: "b7", corner: true, row: 6, col: 1, name: "testname"},
		{in: "b7:b9", corner: false, row: 6, col: 1, name: "testname"},
		{in: "b7:b9", corner: true, row: 8, col: 1, name: "testname"},
		{in: "b7:c9", corner: false, row: 6, col: 1, name: "testname"},
		{in: "b7:c9", corner: true, row: 8, col: 2, name: "testname"},
		{in: "b7:c", corner: false, row: 6, col: 1, name: "testname"},
		{in: "b7:c", corner: true, row: 9999, col: 2, name: "testname"},
		{in: "b7:", corner: false, row: 6, col: 1, name: "testname"},
		{in: "b7:", corner: true, row: 9999, col: 9999, name: "testname"},
	}
	_ = tbl
	for _, testcase := range tbl {
		rowRes, colRes := ssdb.ParseSymbolicRange(testcase.corner, testcase.in)
		assert.Equal(t, testcase.row, rowRes, testcase.name+" "+testcase.in)
		assert.Equal(t, testcase.col, colRes, testcase.name+testcase.in)
	}
}

func TestParseSymbolicCell(t *testing.T) {
	tbl := []struct {
		in     string
		corner bool // false is NW, true is SE
		row    int
		col    int
	}{
		{in: "b7", corner: false, row: 6, col: 1},
		{in: "b7", corner: true, row: 6, col: 1},
		{in: "b", corner: false, row: 0, col: 1},
		{in: "b", corner: true, row: 9999, col: 1},
		{in: "7", corner: false, row: 6, col: 0},
		{in: "7", corner: true, row: 6, col: 9999},
	}
	_ = tbl
	for _, testcase := range tbl {
		rowRes, colRes := ssdb.ParseSymbolicCell(testcase.corner, testcase.in)
		assert.Equal(t, testcase.row, rowRes)
		assert.Equal(t, testcase.col, colRes)
	}
}

func TestParseSymbolicColumn(t *testing.T) {
	tbl := []struct {
		in  string
		out int
	}{
		{in: "A", out: 0},
		{in: "B", out: 1},
		{in: "AA", out: 26},
		{in: "AB", out: 27},
		{in: "BA", out: 52},
		{in: "AAA", out: -1},
		{in: "", out: -1},
	}
	_ = tbl
	for _, testcase := range tbl {
		res := ssdb.ParseSymbolicColumn(testcase.in)
		assert.Equal(t, testcase.out, res)
	}
}
