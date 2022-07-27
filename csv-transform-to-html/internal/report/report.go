package report

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/dele454/medium/csv-transform-to-html/internal/errs"
	"github.com/dele454/medium/csv-transform-to-html/internal/utils"
)

// Reporter ops for any tranformation reporter
type Reporter interface {
	AddError(err error)
	RecordFailed()
	RecordProcessed()
	RecordTransformed()
	Completed()
	AddDuration(since float64)

	WriteReportToFile(ctx context.Context) error
	WriteReportToStdOut(ctx context.Context) error

	SetFilename(name string)
	SetHeaders(headers []string)

	GetHeaders() []string
	GetFilename() string
	GetErrors() []error
	GetTotalTransformedRecords() int
	GetTotalFailedRecords() int
}

// TransformationReporter provides stats after the complete
// processing of a credit source file
type TransformationReporter struct {
	FileName                string
	Headers                 []string
	SendReportToFile        bool
	TotalProcessedRecords   int
	TotalTransformedRecords int
	TotalFailedRecords      int
	Errors                  []error
	Duration                float64
	DurationDisplay         string
	CompletedAt             string
}

// NewTransformationReporter create a new instance of a report
//
// Accepts a flag on whether to send report to a file or not
// If false stdout is used as the ouput.
func NewTransformationReporter(s bool) *TransformationReporter {
	return &TransformationReporter{
		SendReportToFile: s,
	}
}

// WriteReportToFile writes report to file
func (t *TransformationReporter) WriteReportToFile(ctx context.Context) error {
	var err error

	if !t.SendReportToFile {
		return t.WriteReportToStdOut(ctx)
	}

	if t.FileName == "" {
		return errs.ErrorNoSourceFileName
	}

	path := utils.RootDir()

	// create template
	tmpl, err := template.Must(template.New("TEXT"), err).ParseFiles(path + "/internal/transform/template/report.tmpl")
	if err != nil {
		return err
	}

	// apply tmpl to data
	var processed bytes.Buffer
	err = tmpl.ExecuteTemplate(&processed, "report.tmpl", t)
	if err != nil {
		return err
	}

	// create report folder if not exists
	folder := path + "/report"

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err := os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return fmt.Errorf(errs.ErrorFailedToCreateDirectory.Error(), err)
		}
	}

	// create output txt file
	f, err := os.Create(fmt.Sprintf("%s/%s_%s.txt", folder, t.FileName, t.CompletedAt))
	if err != nil {
		utils.Log(utils.ColorError, err)
		return err
	}

	// write output to txt file
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

// WriteReportToStdOut writes report to stdout
func (t *TransformationReporter) WriteReportToStdOut(ctx context.Context) error {
	var err error

	path := utils.RootDir()

	// create template
	tmpl, err := template.Must(template.New("STDOUT"), err).
		ParseFiles(path + "/internal/transform/template/report.tmpl")
	if err != nil {
		return err
	}

	// apply tmpl to data
	var processed bytes.Buffer
	err = tmpl.ExecuteTemplate(&processed, "report.tmpl", t)
	if err != nil {
		return err
	}

	// write output to stdout
	w := bufio.NewWriter(os.Stdout)
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

// AddError adds any errors found to error report
func (t *TransformationReporter) AddError(err error) {
	t.Errors = append(t.Errors, err)
}

// RecordFailed increments the failed records count
func (t *TransformationReporter) RecordFailed() {
	t.TotalFailedRecords++
}

// RecordTransformed increments the transformed records count
func (t *TransformationReporter) RecordTransformed() {
	t.TotalTransformedRecords++
}

// RecordProcessed increments the processed records count
func (t *TransformationReporter) RecordProcessed() {
	t.TotalProcessedRecords++
}

// Completed records the ts of the entire process
func (t *TransformationReporter) Completed() {
	t.CompletedAt = time.Now().Format(time.RFC3339)
	t.DurationDisplay = fmt.Sprintf("%.2f", t.Duration) + "s"
}

// AddDuration increments the processed records count
func (t *TransformationReporter) AddDuration(since float64) {
	t.Duration += since
}

// SetFilename sets the name of the file
func (t *TransformationReporter) SetFilename(name string) {
	t.FileName = name
}

// GetFilename gets the name of the file
func (t *TransformationReporter) GetFilename() string {
	return t.FileName
}

// SetHeaders sets the name of the file
func (t *TransformationReporter) SetHeaders(headers []string) {
	t.Headers = headers
}

// GetHeaders gets the headers of the file
func (t *TransformationReporter) GetHeaders() []string {
	return t.Headers
}

// GetErrors returns all the errors in the report
func (t *TransformationReporter) GetErrors() []error {
	return t.Errors
}

// GetTotalTransformedRecords returns the total transformed records
func (t *TransformationReporter) GetTotalTransformedRecords() int {
	return t.TotalTransformedRecords
}

// GetTotalFailedRecords returns the total failed records
func (t *TransformationReporter) GetTotalFailedRecords() int {
	return t.TotalFailedRecords
}
