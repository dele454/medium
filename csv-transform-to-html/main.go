package main

import (
	"flag"
	"os"

	"github.com/dele454/medium/csv-transform-to-html/cmd"
	"github.com/dele454/medium/csv-transform-to-html/internal/errs"
	"github.com/dele454/medium/csv-transform-to-html/internal/utils"
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
		utils.Log(utils.ColorError, err)
		return
	}

	// check its a file
	if f.IsDir() {
		utils.Log(utils.ColorError, errs.ErrorArgsDirSpecified)
		return
	}

	// kickoff the process
	cmd.Process(file)
}
