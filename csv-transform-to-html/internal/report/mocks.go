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
	TotalProcessedRecords   int
	TotalTransformedRecords int
	TotalFailedRecords      int
	Errors                  []error
	Duration                float64
	DurationDisplay         string
	CompletedAt             string
}

// NewMockReporter create a new instance of a mock reporter
func NewMockReporter() *Mock {
	return &Mock{}
}

// WriteReportToStdOut writes report to stdout
func (m *Mock) WriteReportToStdOut(ctx context.Context) error {
	return nil
}

// AddError adds any errors found to error report
func (m *Mock) AddError(err error) {
	m.Errors = append(m.Errors, err)
}

// RecordFailed increments the failed records count
func (m *Mock) RecordFailed() {
	m.TotalFailedRecords++
}

// RecordTransformed increments the transformed records count
func (m *Mock) RecordTransformed() {
	m.TotalTransformedRecords++
}

// RecordProcessed increments the processed records count
func (m *Mock) RecordProcessed() {
	m.TotalProcessedRecords++
}

// Completed records the ts of the entire process
func (m *Mock) Completed() {
	m.CompletedAt = time.Now().Format(time.RFC3339)
	m.DurationDisplay = fmt.Sprintf("%.2f", m.Duration) + "s"
}

// AddDuration increments the processed records count
func (m *Mock) AddDuration(since float64) {
	m.Duration += since
}

// GetErrors returns all the errors in the report
func (m *Mock) GetErrors() []error {
	return m.Errors
}

// GetTotalTransformedRecords returns the total transformed records
func (m *Mock) GetTotalTransformedRecords() int {
	return m.TotalTransformedRecords
}

// GetTotalFailedRecords returns the total failed records
func (m *Mock) GetTotalFailedRecords() int {
	return m.TotalFailedRecords
}

// SetFilename sets the name of the file
func (m *Mock) SetFilename(name string) {
	m.FileName = name
}

// GetFilename gets the name of the file
func (m *Mock) GetFilename() string {
	return m.FileName
}

// GetHeaders gets the headers of the file
func (m *Mock) GetHeaders() []string {
	return m.Headers
}

// SetHeaders sets the name of the file
func (m *Mock) SetHeaders(headers []string) {
	m.Headers = headers
}
