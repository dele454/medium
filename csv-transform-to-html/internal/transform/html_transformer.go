package transform

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"time"

	"github.com/dele454/medium/csv-transform-to-html/internal/errs"
	"github.com/dele454/medium/csv-transform-to-html/internal/report"
	"github.com/dele454/medium/csv-transform-to-html/internal/utils"
)

// HTMLTransformer creates a new instance of a transformer
//
// Accepts a reporter for reporting purposes.
func NewHTMLTransformer(reporter report.Reporter) *HTMLTransformer {
	return &HTMLTransformer{
		reporter: reporter,
	}
}

// ProcessRecord process records received via the chan
func (tr *HTMLTransformer) ProcessRecord(wg *sync.WaitGroup, record <-chan []string, done <-chan bool) {
	var (
		now  = time.Now()
		data [][]string
		end  bool
		ctx  = context.Background()
	)

	defer func() {
		tr.reporter.SetFilename(filepath.Base(tr.reporter.GetFilename()))
		tr.reporter.AddDuration(time.Since(now).Seconds())
		tr.reporter.Completed()

		err := tr.reporter.WriteReportToFile(ctx)
		if err != nil {
			utils.Log(utils.ColorError, err)
		}

		utils.Log(utils.ColorOK, fmt.Sprintf("Done transforming and exporting %s entries.",
			filepath.Ext(tr.reporter.GetFilename())))

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

			// push row to collection
			tr.reporter.RecordTransformed()
			data = append(data, row)
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
