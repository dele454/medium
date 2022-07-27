package report

import (
	"context"
	"fmt"
	"time"
)

// Mock is a mock reporter for testing
type Mock struct {
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

// NewMockReporter create a new instance of a mock reporter
func NewMockReporter(s bool) *Mock {
	return &Mock{
		SendReportToFile: s,
	}
}

// WriteReportToFile writes report to file
func (t *Mock) WriteReportToFile(ctx context.Context) error {
	return nil
}

// WriteReportToStdOut writes report to stdout
func (t *Mock) WriteReportToStdOut(ctx context.Context) error {
	return nil
}

// AddError adds any errors found to error report
func (t *Mock) AddError(err error) {
	t.Errors = append(t.Errors, err)
}

// RecordFailed increments the failed records count
func (t *Mock) RecordFailed() {
	t.TotalFailedRecords++
}

// RecordTransformed increments the transformed records count
func (t *Mock) RecordTransformed() {
	t.TotalTransformedRecords++
}

// RecordProcessed increments the processed records count
func (t *Mock) RecordProcessed() {
	t.TotalProcessedRecords++
}

// Completed records the ts of the entire process
func (t *Mock) Completed() {
	t.CompletedAt = time.Now().Format(time.RFC3339)
	t.DurationDisplay = fmt.Sprintf("%.2f", t.Duration) + "s"
}

// AddDuration increments the processed records count
func (t *Mock) AddDuration(since float64) {
	t.Duration += since
}

// GetErrors returns all the errors in the report
func (t *Mock) GetErrors() []error {
	return t.Errors
}

// GetTotalTransformedRecords returns the total transformed records
func (t *Mock) GetTotalTransformedRecords() int {
	return t.TotalTransformedRecords
}

// GetTotalFailedRecords returns the total failed records
func (t *Mock) GetTotalFailedRecords() int {
	return t.TotalFailedRecords
}

// SetFilename sets the name of the file
func (t *Mock) SetFilename(name string) {
	t.FileName = name
}

// GetFilename gets the name of the file
func (t *Mock) GetFilename() string {
	return t.FileName
}

// GetHeaders gets the headers of the file
func (t *Mock) GetHeaders() []string {
	return t.Headers
}

// SetHeaders sets the name of the file
func (t *Mock) SetHeaders(headers []string) {
	t.Headers = headers
}
