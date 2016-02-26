// Copyright (c) 2016, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD-style license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/js-arias/cmdapp"
)

var colsCmd = &cmdapp.Command{
	Run: colsRun,
	UsageLine: `cols [-f <char>] [-i|--input <file>] [-n|--no-header]
	[-o|--output <file>] [-v|--invert] <column>...`,
	Short: "selects columns by name",
	Long: `
Command cols selects columns by name and outputs a table with that columns.
If a column name does not match any of the columns in the table, cols creates
a new empty column by that name in the indicated location.

If no columns are indicated all the columns in the table will be selected.

This command can be used to select, sort, or delete columns.

Options are:

    -f <char>
      Sets the field separation character. By default the value is the tab
      character.

    -i <file>
    --input <file>
      Read the table from <file> instead of stdin.

    -n
    --no-header
      If set, the table will be printed without a header.

    -o <file>
    --output <file>
      Write the resulting table to <file> instead of stdout. 

    -v
    --invert
      Inverts the program behavior, i.e. output only the columns NOT included
      in the arguments.

    <column>
      One or more column names.
	`,
}

func init() {
	initCommonFlags(colsCmd)
}

func colsRun(c *cmdapp.Command, args []string) error {
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
	var head []int
	var cols []string
	if invert {
		var err error
		cols, head, err = deleteColumns(r, args)
		if err != nil {
			return err
		}
	} else {
		var err error
		cols, head, err = selectColumns(r, args)
		if err != nil {
			return err
		}
	}

	if len(head) == 0 {
		return nil
	}
	if !noHead {
		err := w.Write(cols)
		if err != nil {
			return err
		}
	}

	for {
		row, err := colsFn(r, head)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		err = w.Write(row)
		if err != nil {
			return err
		}
	}
	return nil
}

// colsFn return a row with the columns indicated by the head slice.
func colsFn(r *csv.Reader, head []int) (row []string, err error) {
	nr, err := r.Read()
	if err != nil {
		return nil, err
	}
	row = make([]string, len(head))
	for i, h := range head {
		if h == -1 {
			continue
		}
		row[i] = nr[h]
	}
	return row, nil
}

// selectColumns selects the columns indicated by args, it returns an slice
// with the column names of the new table, and an int slice with the column
// order (-1 if the column is new) on the original table. If no columns are
// indicated, it will return all the columns in the original table.
func selectColumns(r *csv.Reader, args []string) (cols []string, head []int, err error) {
	header, err := r.Read()
	if err != nil {
		return nil, nil, err
	}

	// if no columns are given returns all columns
	if len(args) == 0 {
		for i := range header {
			head = append(head, i)
		}
		cols = header
		return
	}

	// lookup for columns
	head = make([]int, len(args))
	cols = make([]string, len(args))
	copy(cols, args)
	for i, c := range args {
		head[i] = -1
		for j, h := range header {
			if c == h {
				head[i] = j
				break
			}
		}
	}
	return
}

// deleteColumns returns a slice with columns names of the new table, and an
// int slice with the number of the retained columns in the original table.
func deleteColumns(r *csv.Reader, args []string) (cols []string, head []int, err error) {
	header, err := r.Read()
	if err != nil {
		return nil, nil, err
	}

	// if no column are given returns an empty head index
	if len(args) == 0 {
		return nil, nil, nil
	}

	// removes the columns found
	for i, h := range header {
		toDel := false
		for _, c := range args {
			if c == h {
				toDel = true
				break
			}
		}
		if toDel {
			continue
		}
		head = append(head, i)
		cols = append(cols, h)
	}
	return
}
