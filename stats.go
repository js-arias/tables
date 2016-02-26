// Copyright (c) 2016, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD-style license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/js-arias/cmdapp"
)

var statsCmd = &cmdapp.Command{
	Run: statsRun,
	UsageLine: `stats [-f <char>] [-i|--input <file>] [-o|--output <file>]
	[-p <number>] [-z|--empty-as-zero] <column>...`,
	Short: "calculate basic stats of columns",
	Long: `
Command stats reads an input table and prints on the standard output a new
table with the basic statistics of the indicated columns.

By default empty or non-numeric columns in a row will be ignored, if option
-z or --empty-as-zero is used, that columns will be interpreted as having a
zero value.

Options are:

    -f <char>
      Sets the field separation charachter. By default the value is the tab
      character.

    -i <file>
    --input <file>
      Read the table from <file> instead of stdin.

    -o <file>
    --output <file>
      Write the resulting table to <file> instead of stdout. 

    -p <number>
      Sets the precision in number of decimals. The default is 3.

    -z
    --empty-as-zero
      If set, empty cells, or cells with non-numeric values will be counted as
      having a zero. Otherwise, they will be ignored.

    <column>
      One or more column names.
	`,
}

var emptyZero bool // set empty fields as zero, -z|--empty-as-zero
var precVal int    // set precission, -p

func init() {
	initCommonFlags(statsCmd)
	statsCmd.Flag.BoolVar(&emptyZero, "empty-as-zero", false, "")
	statsCmd.Flag.BoolVar(&emptyZero, "z", false, "")
	statsCmd.Flag.IntVar(&precVal, "p", 3, "")
}

func statsRun(c *cmdapp.Command, args []string) error {
	in := os.Stdin
	if len(input) > 0 {
		var err error
		in, err = os.Open(input)
		if err != nil {
			return err
		}
		defer in.Close()
	}
	out := os.Stdout
	if len(output) > 0 {
		var err error
		out, err = os.Create(output)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	if len(delim) == 0 {
		delim = "\t"
	}
	r1 := []rune(delim)[0]
	r := csv.NewReader(in)
	r.Comma = r1
	w := csv.NewWriter(out)
	w.Comma = r1
	w.UseCRLF = true
	defer w.Flush()
	cols, head, err := selectColumns(r, args)
	if err != nil {
		return err
	}
	calc := make([]statsCalc, len(head))
	for {
		row, oks, err := statsFn(r, head)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		for i := range calc {
			if (!oks[i]) && (!emptyZero) {
				continue
			}
			calc[i].sum += row[i]
			if calc[i].n == 0 {
				calc[i].max = row[i]
				calc[i].min = row[i]
			}
			calc[i].n++
			if calc[i].max < row[i] {
				calc[i].max = row[i]
			}
			if calc[i].min > row[i] {
				calc[i].min = row[i]
			}
			prev := calc[i].a
			calc[i].a = calc[i].a + ((row[i] - calc[i].a) / float64(calc[i].n))
			calc[i].q = calc[i].q + ((row[i] - prev) * (row[i] - calc[i].a))
		}
	}

	// writes output header
	header := []string{"Stat"}
	header = append(header, cols...)
	err = w.Write(header)
	if err != nil {
		return err
	}

	// writes sum
	row := []string{"Sum"}
	for i := range calc {
		row = append(row, strconv.FormatFloat(calc[i].sum, 'g', precVal, 64))
	}
	err = w.Write(row)
	if err != nil {
		return err
	}

	// writes mean
	row = []string{"Mean"}
	for i := range calc {
		if calc[i].n == 0 {
			row = append(row, "NaN")
			continue
		}
		row = append(row, strconv.FormatFloat(calc[i].a, 'g', precVal, 64))
	}
	err = w.Write(row)
	if err != nil {
		return err
	}

	// writes max
	row = []string{"Max"}
	for i := range calc {
		if calc[i].n == 0 {
			row = append(row, "NaN")
			continue
		}
		row = append(row, strconv.FormatFloat(calc[i].max, 'g', precVal, 64))
	}
	err = w.Write(row)
	if err != nil {
		return err
	}

	// writes min
	row = []string{"Min"}
	for i := range calc {
		if calc[i].n == 0 {
			row = append(row, "NaN")
			continue
		}
		row = append(row, strconv.FormatFloat(calc[i].min, 'g', precVal, 64))
	}
	err = w.Write(row)
	if err != nil {
		return err
	}

	// writes standard deviation
	row = []string{"StDev"}
	for i := range calc {
		if calc[i].n < 1 {
			row = append(row, "NaN")
			continue
		}
		s2 := calc[i].q / float64(calc[i].n-1)
		row = append(row, strconv.FormatFloat(math.Sqrt(s2), 'g', precVal, 64))
	}
	err = w.Write(row)
	if err != nil {
		return err
	}

	// writes range
	row = []string{"Range"}
	for i := range calc {
		if calc[i].n == 0 {
			row = append(row, "NaN")
			continue
		}
		row = append(row, strconv.FormatFloat(calc[i].max-calc[i].min, 'g', precVal, 64))
	}
	err = w.Write(row)
	if err != nil {
		return err
	}

	return nil
}

// statsCalc contains the variables to calculate basic stats.
type statsCalc struct {
	sum float64
	n   int
	max float64
	min float64
	a   float64 // mean
	q   float64 // variance sum
}

// statsFn returns the numeric values of a set of columns (defined by head) in
// a table. If no value is found, a zero will be returned and the
// corresponding ok value as false.
func statsFn(r *csv.Reader, head []int) (row []float64, oks []bool, err error) {
	nr, err := r.Read()
	if err != nil {
		return nil, nil, err
	}
	row = make([]float64, len(head))
	oks = make([]bool, len(head))
	for i, h := range head {
		if h == -1 {
			continue
		}
		v, err := strconv.ParseFloat(nr[h], 64)
		if err != nil {
			continue
		}
		row[i] = v
		oks[i] = true
	}
	return
}
