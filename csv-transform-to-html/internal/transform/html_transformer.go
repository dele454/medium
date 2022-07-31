package transform

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dele454/medium/csv-transform-to-html/internal/errs"
	"github.com/dele454/medium/csv-transform-to-html/internal/report"
	"github.com/dele454/medium/csv-transform-to-html/internal/utils"
)

// Output details of the transformation sent to an io.Writer
type Output struct {
	FileName     string
	FileLocation string
	TotalHeaders int
	Headers      []string
	Data         []utils.SalesRecord
}

// Transformer ops every transformer should conform to
type Transformer interface {
	WriteOutputToFile(output *Output) error
	ProcessRecord(wg *sync.WaitGroup, record <-chan []string, done <-chan bool)
}

// HTMLTransformer handles the processing of customer data
// with the aid of preprocessors and writing the output
// as an HTML document.
type HTMLTransformer struct {
	processor utils.PreProcessor
	reporter  report.Reporter
}

// XMLTransformer handles the processing of customer data
// with the aid of preprocessors and writing the output
// as an XML document.
type XMLTransformer struct{}

// HTMLTransformer creates a new instance of a transformer
//
// Accepts a reporter for reporting purposes.
func NewHTMLTransformer(reporter report.Reporter) Transformer {
	return &HTMLTransformer{
		processor: utils.NewProcessor(),
		reporter:  reporter,
	}
}

// ProcessRecord process records received via the chan
func (tr *HTMLTransformer) ProcessRecord(wg *sync.WaitGroup, record <-chan []string, done <-chan bool) {
	var (
		now  = time.Now()
		data []utils.SalesRecord
		end  bool
		ctx  = context.Background()
	)

	defer func() {
		tr.reporter.SetFilename(filepath.Base(tr.reporter.GetFilename()))
		tr.reporter.AddDuration(time.Since(now).Seconds())
		tr.reporter.Completed()

		err := tr.reporter.WriteReportToStdOut(ctx)
		if err != nil {
			utils.Log(utils.ColorError, err)
		}

		wg.Done()
	}()

	// process pipeline
	for {
		select {
		case <-done:
			// reading has completed.
			end = true
		case row := <-record:
			// read from pipeline
			if len(row) == 0 {
				tr.reporter.RecordFailed()
				tr.reporter.AddError(errs.ErrorEmptyRowFound)

				utils.Log(utils.ColorError, errs.ErrorEmptyRowFound)
				continue
			}

			// unmarshal records
			var sr utils.SalesRecord
			sr, err := tr.processor.Unmarshal(row, sr)
			if err != nil {
				tr.reporter.RecordFailed()
				tr.reporter.AddError(err)

				continue
			}

			// push row to collection
			tr.reporter.RecordTransformed()
			data = append(data, sr)
		}

		// means parser has signaled end of file
		// exit loop
		if end {
			break
		}
	}

	// send output to file
	err := tr.WriteOutputToFile(&Output{
		FileName:     filepath.Base(tr.reporter.GetFilename()),
		TotalHeaders: len(tr.reporter.GetHeaders()),
		Headers:      tr.reporter.GetHeaders(),
		Data:         data,
	})
	if err != nil {
		tr.reporter.AddError(err)
	}
}

// WriteOutputToFile write output data to file.
func (tr *HTMLTransformer) WriteOutputToFile(output *Output) error {
	var err error

	path := utils.RootDir()

	// create template
	tmpl, err := template.Must(template.New("HTML"), err).ParseFiles(path + "/internal/transform/template/output.tmpl")
	if err != nil {
		tr.reporter.AddError(err)
		return err
	}

	// apply tmpl to data
	var processed bytes.Buffer
	err = tmpl.ExecuteTemplate(&processed, "output.tmpl", output)
	if err != nil {
		return err
	}

	// create output folder if not exists
	folder := path + "/output"

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err := os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return fmt.Errorf(errs.ErrorFailedToCreateDirectory.Error(), err)
		}
	}

	// create output html file
	f, err := os.Create(fmt.Sprintf("%s/%s.html", folder, output.FileName))
	if err != nil {
		return err
	}

	// write output to html file
	w := bufio.NewWriter(f)
	_, err = w.WriteString(processed.String())
	if err != nil {
		return err
	}

	// flush buffer
	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}
