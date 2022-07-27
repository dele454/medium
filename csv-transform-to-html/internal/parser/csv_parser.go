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
	GetFileInfo() []string
	Read(wg *sync.WaitGroup)
	Close() error
}

// CSVParser parser for parsing and reading from csv files
type CSVParser struct {
	// file   CSVFile
	report report.Reporter
}

// CSVFile
type CSVFile struct {
	Name    string
	Headers []string
}

// NewCSVParser creates csv parser for parsing & reading credit info from a csv file
func NewCSVParser(file string, reporter report.Reporter) *CSVParser {
	reporter.SetFilename(file)

	return &CSVParser{
		report: reporter,
	}
}

// Read reads from the source file
func (c *CSVParser) Read(wg *sync.WaitGroup, record chan<- []string, done chan<- bool) {
	start := time.Now()

	defer func() {
		close(done)
		close(record)
		c.report.AddDuration(time.Since(start).Seconds())

		utils.Log(utils.ColorOK, "Done reading & parsing csv file.")

		wg.Done()
	}()

	f, err := os.Open(c.report.GetFilename())
	if err != nil {
		utils.Log(utils.ColorError, fmt.Sprintf("Error opening csv file: %s", err))
		return
	}
	defer f.Close()

	// create reader
	reader := csv.NewReader(f)

	// read the headers
	headers, err := reader.Read()
	if err != nil {
		utils.Log(utils.ColorError, errs.ErrorNoHeadersFound)
		return
	}

	// set the headers
	c.report.SetHeaders(headers)
	c.report.SetFilename(filepath.Base(c.report.GetFilename()))

	// read from file
	for {
		row, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				c.report.AddError(err)
				c.report.RecordFailed()
				continue
			}
			break
		}

		c.report.RecordProcessed()
		record <- row
	}

	done <- true
}
