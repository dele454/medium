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
	reporter report.Reporter
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
		reporter: reporter,
	}
}

// Read reads from the source file
func (c *CSVParser) Read(wg *sync.WaitGroup, record chan<- []string, done chan<- bool) {
	start := time.Now()

	defer func() {
		close(done)
		close(record)
		c.reporter.AddDuration(time.Since(start).Seconds())

		wg.Done()
	}()

	f, err := os.Open(c.reporter.GetFilename())
	if err != nil {
		c.reporter.AddError(err)
		utils.Log(utils.ColorError, fmt.Sprintf("Error opening csv file: %s", err))
		return
	}
	defer f.Close()

	// create reader
	reader := csv.NewReader(f)

	// read the headers
	headers, err := reader.Read()
	if err != nil {
		c.reporter.AddError(err)
		utils.Log(utils.ColorError, errs.ErrorNoHeadersFound)
		return
	}

	// set the headers
	c.reporter.SetHeaders(headers)
	c.reporter.SetFilename(filepath.Base(c.reporter.GetFilename()))

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
