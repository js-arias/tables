// Copyright (c) 2016, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD-style license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/js-arias/cmdapp"
)

var rowsCmd = &cmdapp.Command{
	Run: rowsRun,
	UsageLine: `rows [-f <char>] [-i|--input <file>] [-n|--no-header]
	[-o|--output <file>] [-v|--invert] <expression>...`,
	Short: "Select rows matching an expression",
	Long: `
Command rows select rows that fullfill the conditions given in the expression.
An expression must start with a column name followed by a conditional operand
("==", "!=", "<", "<=", ">", ">=") and the a column name, an string bounded
by quotes ("), or a number.

When multiple expressions are indicated they are taken as an or codition, to
implement and, pipe this command.

This command is meant to be used interactively at the command-line, so it
makes relatively simple operations. Programs or scripts should be preferred
for more complex operations.

Because logical conditions are represented by special characters they must be
enclosed in single quotes (') to protect the condition from being interpreted
by the shell.

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

    <expression>
      A conditional expression to be evaluated by rows command.
	`,
}

func init() {
	initCommonFlags(rowsCmd)
}

func rowsRun(c *cmdapp.Command, args []string) error {
	if len(args) == 0 {
		c.Usage()
	}
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
	header, err := r.Read()
	if err != nil {
		return err
	}
	var exps []expression
	for _, a := range args {
		e, err := parseExpression(header, strings.NewReader(strings.TrimSpace(a)))
		if err != nil {
			return err
		}
		exps = append(exps, e)
	}
	w := csv.NewWriter(out)
	w.Comma = r1
	w.UseCRLF = true
	defer w.Flush()
	if !noHead {
		err = w.Write(header)
		if err != nil {
			return err
		}
	}

	for {
		row, err := rowsFn(r, exps)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if len(row) == 0 {
			continue
		}
		err = w.Write(row)
		if err != nil {
			return err
		}
	}
	return nil
}

// rowsFn returns a row if it fullfills the indicated expressions, otherwise
// it returns an empty row.
func rowsFn(r *csv.Reader, exps []expression) (row []string, err error) {
	row, err = r.Read()
	if err != nil {
		return nil, err
	}
	sel := false
	for _, e := range exps {
		val1 := getFieldValue(row[e.cols[0]])
		if val1 == nil {
			continue
		}
		val2 := e.value
		if e.cols[1] != -1 {
			val2 = getFieldValue(row[e.cols[1]])
		}
		if compare(val1, val2, e.op) {
			sel = true
			break
		}
	}
	if (invert && (!sel)) || (sel && (!invert)) {
		return
	}
	return nil, nil
}

// expression defines a comparative expression between values of one or more
// columns.
type expression struct {
	cols  [2]int // column in the expression
	op    int    // operation
	value interface{}
}

// compare returns true if two values fullfills the conditional operator given
// by op.
func compare(val1, val2 interface{}, op int) bool {
	switch v := val1.(type) {
	case string:
		w, ok := val2.(string)
		if !ok {
			// val2 is not an string!
			switch op {
			case opDiff:
				return true
			case opLess:
				// strings are "smaller" than numbers
				return true
			}
			return false
		}
		switch op {
		case opEqual:
			return v == w
		case opDiff:
			return v != w
		case opGreat:
			return v > w
		case opGreatEqual:
			return (v > w) || (v == w)
		case opLess:
			return v < w
		case opLessEqual:
			return (v < w) || (v == w)
		default:
			return false
		}
	case float64:
		w, ok := val2.(float64)
		if !ok {
			// val2 is not a number!
			switch op {
			case opDiff:
				return true
			case opGreat:
				// numbers are "greater" than strings, etc.
				return true
			}
			return false
		}
		x := float64(v)
		switch op {
		case opEqual:
			return x == w
		case opDiff:
			return x != w
		case opGreat:
			return x > w
		case opGreatEqual:
			return x >= w
		case opLess:
			return x < w
		case opLessEqual:
			return x <= w
		default:
			return false
		}
	}
	return false
}

func skipExpressionSpaces(r *strings.Reader) error {
	for {
		r1, _, err := r.ReadRune()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(r1) {
			r.UnreadRune()
			return nil
		}
	}
}

// comparative operators
const (
	opEqual      = iota // ==
	opDiff              // !=
	opGreat             // >
	opGreatEqual        // >=
	opLess              // <
	opLessEqual         // <=
)

// parsExpression returns an expression from a string containing a simple
// comparative expression.
func parseExpression(header []string, r *strings.Reader) (e expression, err error) {
	var b bytes.Buffer

	// get the column name
	for {
		r1, _, err := r.ReadRune()
		if err != nil {
			return expression{}, err
		}
		if unicode.IsSpace(r1) {
			err = skipExpressionSpaces(r)
			if err != nil {
				return expression{}, err
			}
			break
		}
		if (r1 == '=') || (r1 == '!') || (r1 == '>') || (r1 == '<') {
			r.UnreadRune()
			break
		}
		b.WriteRune(r1)
	}
	col1 := b.String()

	// get the operand
	b.Reset()
	for {
		r1, _, err := r.ReadRune()
		if err != nil {
			return expression{}, err
		}
		if unicode.IsSpace(r1) {
			err = skipExpressionSpaces(r)
			if err != nil {
				return expression{}, err
			}
			break
		}
		if r1 == '"' {
			r.UnreadRune()
			break
		}
		if unicode.IsDigit(r1) || (r1 == '-') || (r1 == '.') {
			r.UnreadRune()
			break
		}
		if (r1 != '=') && (r1 != '!') && (r1 != '>') && (r1 != '<') {
			r.UnreadRune()
			break
		}
		b.WriteRune(r1)
	}
	var op int
	switch s := b.String(); s {
	case "==":
		op = opEqual
	case "!=":
		op = opDiff
	case ">":
		op = opGreat
	case ">=":
		op = opGreatEqual
	case "<":
		op = opLess
	case "<=":
		op = opLessEqual
	default:
		return expression{}, fmt.Errorf("unknown operand: %s", s)
	}

	// reads the second operand
	b.Reset()
	r1, _, err := r.ReadRune()
	if err != nil {
		return expression{}, err
	}
	var val interface{}
	var col2 string
	if r1 == '"' {
		// if it has quotes, it is a string
		for {
			r1, _, err = r.ReadRune()
			if err != nil {
				if err == io.EOF {
					break
				}
				return expression{}, err
			}
			if r1 == '"' {
				break
			}
			b.WriteRune(r1)
		}
		val = b.String()
	} else if unicode.IsDigit(r1) || (r1 == '-') || (r1 == '.') {
		// a number
		for {
			b.WriteRune(r1)
			r1, _, err = r.ReadRune()
			if err != nil {
				if err == io.EOF {
					break
				}
				return expression{}, err
			}
			if unicode.IsSpace(r1) {
				break
			}
		}
		val, err = strconv.ParseFloat(b.String(), 64)
		if err != nil {
			return expression{}, err
		}
	} else {
		// another column
		for {
			b.WriteRune(r1)
			r1, _, err = r.ReadRune()
			if err != nil {
				if err == io.EOF {
					break
				}
				return expression{}, err
			}
			if unicode.IsSpace(r1) {
				break
			}
		}
		col2 = b.String()
	}
	e = expression{
		cols:  [2]int{-1, -1},
		op:    op,
		value: val,
	}
	for i, h := range header {
		if h == col1 {
			e.cols[0] = i
		}
		if h == col2 {
			e.cols[1] = i
		}
	}
	if e.cols[0] == -1 {
		return expression{}, errors.New("expecting a valid column name")
	}
	return e, nil
}

// getFieldValue returns the numeric or string value of a row field.
func getFieldValue(field string) (value interface{}) {
	r1 := []rune(field)[0]
	if unicode.IsDigit(r1) || (r1 == '-') || (r1 == '.') {
		var err error
		value, err = strconv.ParseFloat(field, 64)
		if err == nil {
			return
		}
	}
	return field
}
