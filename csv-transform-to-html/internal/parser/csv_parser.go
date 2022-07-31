package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dele454/medium/csv-transform-to-html/internal/errs"
	"github.com/dele454/medium/csv-transform-to-html/internal/report"
	"github.com/dele454/medium/csv-transform-to-html/internal/utils"
)

// Parser list of operations a parser should be able to perform
type Parser interface {
	Read(wg *sync.WaitGroup, record chan<- []string, done chan<- bool)
}

// CSVParser parser for parsing and reading from csv files
type CSVParser struct {
	reporter report.Reporter
}

// CSVFile
type CSVFile struct {
	Name    string
	Headers []string
}

// NewCSVParser creates csv parser for parsing & reading sales data
func NewCSVParser(file string, reporter report.Reporter) Parser {
	reporter.SetFilename(file)

	return &CSVParser{
		reporter: reporter,
	}
}

// Read reads from the csv file
func (c *CSVParser) Read(wg *sync.WaitGroup, record chan<- []string, done chan<- bool) {
	start := time.Now()

	defer func() {
		close(done)
		close(record)
		c.reporter.AddDuration(time.Since(start).Seconds())

		wg.Done()
	}()

	// open file for reading
	f, err := os.Open(c.reporter.GetFilename())
	if err != nil {
		c.reporter.AddError(err)
		utils.Log(utils.ColorError, fmt.Sprintf("Error opening csv file: %s", err))
		return
	}
	defer f.Close()

	// set the headers
	c.reporter.SetHeaders(utils.GetHeaders())

	// set the file name
	c.reporter.SetFilename(filepath.Base(c.reporter.GetFilename()))

	// create csv reader
	reader := csv.NewReader(f)

	// parse headers detected in file
	if err := c.parseHeaders(reader); err != nil {
		c.reporter.AddError(err)
	}

	// read from file
	for {
		row, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				c.reporter.AddError(err)
				c.reporter.RecordFailed()
				continue
			}
			break
		}

		c.reporter.RecordProcessed()
		record <- row
	}

	done <- true
}

func (c *CSVParser) parseHeaders(reader *csv.Reader) error {
	// get headers from file
	headers, err := reader.Read()
	if err != nil {
		return err
	}

	// check for expected nos of headers
	if len(c.reporter.GetHeaders()) != len(headers) {
		return errs.ErrorUnmatachedHeaders
	}

	// check if all expected headers are in source file
	for _, x := range c.reporter.GetHeaders() {
		var found bool

		for _, y := range headers {
			if x == y {
				found = true
			}
		}

		if !found {
			return fmt.Errorf(errs.ErrorHeaderNotFound.Error(), x)
		}
	}

	return nil
}
