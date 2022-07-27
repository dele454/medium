package cmd

import (
	"sync"

	"github.com/dele454/medium/csv-transform-to-html/internal/parser"
	"github.com/dele454/medium/csv-transform-to-html/internal/report"
	"github.com/dele454/medium/csv-transform-to-html/internal/transform"
)

func Process(file string, send bool) {
	// create a reporter
	reporter := report.NewTransformationReporter(send)

	// create a new parser
	parser := parser.NewCSVParser(file, reporter)

	// create waitgroup
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// channels for pipeline
	record := make(chan []string)
	done := make(chan bool)

	// transformer
	go transform.NewHTMLTransformer(reporter).
		ProcessRecord(wg, record, done)

	// source parser
	go parser.Read(wg, record, done)

	// wait for all go routines to finish
	wg.Wait()
}
