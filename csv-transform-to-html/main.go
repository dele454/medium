package main

import (
	"flag"
	"os"

	"github.com/dele454/medium/csv-transform-to-html/cmd"
	"github.com/dele454/medium/csv-transform-to-html/internal/errs"
)

func main() {
	var file string

	// accept arg from stdin
	flag.StringVar(&file, "f", "", "Full path to source file for processing.")
	flag.Parse()

	// display usage if no arg is passed
	if len(os.Args[1:]) == 0 {
		flag.PrintDefaults()
		return
	}

	// get file's info
	f, err := os.Stat(file)
	if err != nil {
		panic(err)
	}

	// check its a file
	if f.IsDir() {
		panic(errs.ErrorArgsDirSpecified)
	}

	// kickoff the process
	cmd.Process(file)
}
