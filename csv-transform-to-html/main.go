package main

import (
	"flag"
	"os"

	"github.com/dele454/medium/csv-transform-to-html/cmd"
	"github.com/dele454/medium/csv-transform-to-html/internal/errs"
	"github.com/dele454/medium/csv-transform-to-html/internal/utils"
)

func main() {
	var (
		file string
		send bool
	)

	// accept params from stdin
	flag.StringVar(&file, "f", "", "Full path to source file for processing.")
	flag.BoolVar(&send, "r", true, "Whether or not to send report to file. Default is true i.e print to file.")
	flag.Parse()

	// display usage if no args are passed
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
	cmd.Process(file, send)
}
