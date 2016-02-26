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

func TestParseExpression(t *testing.T) {
	h := []string{"cost", "number", "id", "name"}

	r := strings.NewReader("cost > 50")
	exp := expression{
		cols:  [2]int{0, -1},
		op:    opGreat,
		value: float64(50),
	}
	testParseExpression(t, h, exp, r)

	r = strings.NewReader("cost<=50")
	exp.op = opLessEqual
	testParseExpression(t, h, exp, r)

	r = strings.NewReader("cost<id")
	exp = expression{
		cols:  [2]int{0, 2},
		op:    opLess,
		value: nil,
	}
	testParseExpression(t, h, exp, r)

	r = strings.NewReader("cost>=id")
	exp.op = opGreatEqual
	testParseExpression(t, h, exp, r)

	r = strings.NewReader(`name == "test name"`)
	exp = expression{
		cols:  [2]int{3, -1},
		op:    opEqual,
		value: "test name",
	}
	testParseExpression(t, h, exp, r)

	r = strings.NewReader(`id!="xABF01"`)
	exp = expression{
		cols:  [2]int{2, -1},
		op:    opDiff,
		value: "xABF01",
	}
	testParseExpression(t, h, exp, r)

}

func testParseExpression(t *testing.T, header []string, exp expression, r *strings.Reader) {
	e, err := parseExpression(header, r)
	if err != nil {
		t.Errorf("Rows: Error while parsing: %v", err)
		return
	}
	if e.cols[0] != exp.cols[0] {
		t.Errorf("Rows: Bad column (1) parsing: expecting %d found %d", exp.cols[0], e.cols[0])
	}
	if e.cols[1] != exp.cols[1] {
		t.Errorf("Rows: Bad column (2) parsing: expecting %d found %d", exp.cols[1], e.cols[1])
	}
	if e.op != exp.op {
		t.Errorf("Rows: Expecting operation %d found %d", exp.op, e.op)
	}
	if e.cols[1] != -1 {
		return
	}
	switch v := exp.value.(type) {
	case float64:
		if f, ok := e.value.(float64); !ok {
			t.Errorf("Rows: Wrong type, expecting %T found %T", v, e.value)
		} else if f != v {
			t.Errorf("Rows: Wrong value, expecting %.3f found %.3f", v, f)
		}
	case string:
		if s, ok := e.value.(string); !ok {
			t.Errorf("Rows: Wrong type, expecting %T found %T", v, e.value)
		} else if s != v {
			t.Errorf("Rows: Wrong value, expecting %s found %s", v, s)
		}
	}
}

func TestCompare(t *testing.T) {
	if !compare("equal", "equal", opEqual) {
		t.Errorf("Rows: \"equal\" == \"equal\" returns false")
	}
	if compare("not-equal", "different", opEqual) {
		t.Errorf("Rows: \"not-equal\" == \"different\" returns true")
	}
	if !compare("not-equal", "different", opDiff) {
		t.Errorf("Rows: \"not-equal\" != \"different\" returns false")
	}
	if compare("equal", "equal", opDiff) {
		t.Errorf("Rows: \"equal\" != \"equal\" returns true")
	}
	if compare("first", "second", opGreat) {
		t.Errorf("Rows: \"first\" > \"second\" returns true")
	}
	if !compare("first", "second", opLess) {
		t.Errorf("Rows: \"first\" < \"second\" returns false")
	}

	if !compare(float64(50), float64(50), opEqual) {
		t.Errorf("Rows: 50 == 50 returns false")
	}
	if compare(float64(50), float64(100), opEqual) {
		t.Errorf("Rows: 50 == 100 returns true")
	}
	if !compare(float64(50), float64(100), opDiff) {
		t.Errorf("Rows: 50 != 100 returns false")
	}
	if compare(float64(50), float64(50), opDiff) {
		t.Errorf("Rows: 50 != 50 returns true")
	}
	if compare(float64(20), float64(100), opGreat) {
		t.Errorf("Rows: 20 > 100 returns true")
	}
	if !compare(float64(20), float64(100), opLess) {
		t.Errorf("Rows: 20 < 100 returns false")
	}

	if compare("50", float64(50), opEqual) {
		t.Errorf("Rows: \"50\" == 50 returns true")
	}
	if !compare("50", float64(50), opDiff) {
		t.Errorf("Rows: \"50\" != 50 returns false")
	}
	if compare("50", float64(50), opGreat) {
		t.Errorf("Rows: \"50\" > 50 returns true")
	}
	if !compare("50", float64(50), opLess) {
		t.Errorf("Rows: \"50\" < 50 returns false")
	}
}

func TestRowsSelect(t *testing.T) {
	// cols blob is in cols_test.go
	r := csv.NewReader(strings.NewReader(colsBlob))
	r.Comma = '\t'
	header, err := r.Read()
	if err != nil {
		t.Errorf("Rows: unexpected error on read: %v", err)
	}
	args := []string{"Cost > 50"}
	var exps []expression
	for _, a := range args {
		e, err := parseExpression(header, strings.NewReader(strings.TrimSpace(a)))
		if err != nil {
			t.Errorf("Rows: unexpected error on expression: %v", err)
		}
		exps = append(exps, e)
	}
	i := 0
	for {
		row, err := rowsFn(r, exps)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("Rows: unexpected error on read: %v", err)
		}
		if len(row) == 0 {
			continue
		}
		i++
	}
	if i != 3 {
		t.Errorf("Rows: expecting %d rows, found: %d", 3, i)
	}
}
