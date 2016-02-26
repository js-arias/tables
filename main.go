// Copyright (c) 2016, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD-style license that can be found in the LICENSE file.

package main

import (
	"github.com/js-arias/cmdapp"
)

func init() {
	cmdapp.Short = "Tables is a tool for management of text based tables."
	cmdapp.Commands = []*cmdapp.Command{
		colsCmd,
		rowsCmd,
		statsCmd,
	}
}

func main() {
	cmdapp.Run()
}

// general flags used by most commands
var (
	delim  string // set field delimitator, -f
	input  string // set input file, -i|--input
	invert bool   // invert command behavior, -v|--invert
	noHead bool   // set the header output, -n|--no-header
	output string // set output file, -o|--output
)

// initialize general flags.
func initCommonFlags(c *cmdapp.Command) {
	c.Flag.StringVar(&delim, "f", "\t", "")
	c.Flag.StringVar(&input, "input", "", "")
	c.Flag.StringVar(&input, "i", "", "")
	c.Flag.BoolVar(&noHead, "no-header", false, "")
	c.Flag.BoolVar(&noHead, "n", false, "")
	c.Flag.StringVar(&output, "output", "", "")
	c.Flag.StringVar(&output, "o", "", "")
	c.Flag.BoolVar(&invert, "invert", false, "")
	c.Flag.BoolVar(&invert, "v", false, "")
}
