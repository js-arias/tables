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

func TestStats(t *testing.T) {
	h := []string{"Cost", "Value", "Description"}
	r := csv.NewReader(strings.NewReader(colsBlob))
	r.Comma = '\t'
	_, head, err := selectColumns(r, h)
	if err != nil {
		t.Errorf("Stats: unexpected error: %v", err)
	}
	for i := 0; ; i++ {
		row, oks, err := statsFn(r, head)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("Stats: unexpected error: %v", err)
		}
		if i != 0 {
			continue
		}
		if row[0] != 50 {
			t.Errorf("Stats: expecting %.3f found %.3f", 50, row[0])
		}
		if oks[2] {
			t.Errorf("Stats: column %d should be false", 2)
		}
	}
}
