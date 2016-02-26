// Copyright (c) 2016, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD-style license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"io"
	"strings"
	"testing"
)

var colsBlob = `
Item	Amount	Cost	Value	Description
1	3	50	150	rubber gloves
2	100	5	500	test tubes
3	5	80	400	clamps
4	23	19	437	plates
5	99	24	2376	cleaning cloth
6	89	147	13083	bunsen burners
7	5	175	875	scales
`

func TestColsSelect(t *testing.T) {
	h := []string{"Item", "Cost", "Amount"}
	x := []int{0, 2, 1}
	r := csv.NewReader(strings.NewReader(colsBlob))
	r.Comma = '\t'
	cols, head, err := selectColumns(r, h)
	if err != nil {
		t.Errorf("Cols: unexpected error: %v", err)
	}
	if len(cols) != len(head) {
		t.Errorf("Cols: length of cols (%d) and head (%d) differnet", len(cols), len(head))
	}
	for i, v := range h {
		if v != cols[i] {
			t.Errorf("Cols: expecting %s found %s", v, cols[i])
		}
		if x[i] != head[i] {
			t.Errorf("Cols: expecting %d index, found %d", x[i], head[i])
		}
	}

	rs := [][]string{
		[]string{"1", "50", "3"},
		[]string{"2", "5", "100"},
		[]string{"3", "80", "5"},
	}
	for i := 0; ; i++ {
		row, err := colsFn(r, head)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("Cols: unexpected error: %v", err)
		}
		if len(row) != len(head) {
			t.Errorf("Cols: expecting vector of %d elements, found %d (row %d)", len(head), len(row))
		}
		if i < len(rs) {
			for j, v := range rs[i] {
				if row[j] != v {
					t.Errorf("Cols: expecting %s in row %d col %d, found %s", v, i, j, row[j])
				}
			}
		}
	}
}

func TestAddCols(t *testing.T) {
	h := []string{"Item", "Cost", "Amount", "Total"}
	x := []int{0, 2, 1, -1}
	r := csv.NewReader(strings.NewReader(colsBlob))
	r.Comma = '\t'
	cols, head, err := selectColumns(r, h)
	if err != nil {
		t.Errorf("Cols: unexpected error: %v", err)
	}
	if len(cols) != len(head) {
		t.Errorf("Cols: length of cols (%d) and head (%d) differnet", len(cols), len(head))
	}
	for i, v := range h {
		if v != cols[i] {
			t.Errorf("Cols: expecting %s found %s", v, cols[i])
		}
		if x[i] != head[i] {
			t.Errorf("Cols: expecting %d index, found %d", x[i], head[i])
		}
	}
	rs := [][]string{
		[]string{"1", "50", "3", ""},
		[]string{"2", "5", "100", ""},
		[]string{"3", "80", "5", ""},
	}
	for i := 0; ; i++ {
		row, err := colsFn(r, head)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("Cols: unexpected error: %v", err)
		}
		if len(row) != len(head) {
			t.Errorf("Cols: expecting vector of %d elements, found %d (row %d)", len(head), len(row))
		}
		if i < len(rs) {
			for j, v := range rs[i] {
				if row[j] != v {
					t.Errorf("Cols: expecting %s in row %d col %d, found %s", v, i, j, row[j])
				}
			}
		}
	}
}

func TestDelCols(t *testing.T) {
	h := []string{"Item", "Cost", "Amount"}
	header := []string{"Value", "Description"}
	x := []int{3, 4}
	r := csv.NewReader(strings.NewReader(colsBlob))
	r.Comma = '\t'
	cols, head, err := deleteColumns(r, h)
	if err != nil {
		t.Errorf("Cols: unexpected error: %v", err)
	}
	if len(cols) != len(head) {
		t.Errorf("Cols: length of cols (%d) and head (%d) differnet", len(cols), len(head))
	}
	if len(cols) != len(header) {
		t.Errorf("Cols: length of cols (%d) and header (%d) differnet", len(cols), len(header))
	}
	for i, v := range header {
		if v != cols[i] {
			t.Errorf("Cols: expecting %s found %s", v, cols[i])
		}
		if x[i] != head[i] {
			t.Errorf("Cols: expecting %d index, found %d", x[i], head[i])
		}
	}
	rs := [][]string{
		[]string{"150", "rubber gloves"},
		[]string{"500", "test tubes"},
		[]string{"400", "clamps"},
		[]string{"437", "plates"},
	}
	for i := 0; ; i++ {
		row, err := colsFn(r, head)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("Cols: unexpected error: %v", err)
		}
		if len(row) != len(head) {
			t.Errorf("Cols: expecting vector of %d elements, found %d (row %d)", len(head), len(row))
		}
		if i < len(rs) {
			for j, v := range rs[i] {
				if row[j] != v {
					t.Errorf("Cols: expecting %s in row %d col %d, found %s", v, i, j, row[j])
				}
			}
		}
	}
}
