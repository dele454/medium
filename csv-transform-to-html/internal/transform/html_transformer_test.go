package transform

import (
	"sync"
	"testing"

	"github.com/dele454/medium/csv-transform-to-html/internal/parser"
	"github.com/dele454/medium/csv-transform-to-html/internal/report"
	"github.com/dele454/medium/csv-transform-to-html/internal/utils"
)

func TestProcessRecord(t *testing.T) {
	reporter := report.NewMockReporter(false)
	transformer := NewHTMLTransformer(reporter)

	path := utils.RootDir()
	p := parser.NewCSVParser(path+"/internal/testdata/100_sales_records.csv", reporter)

	// create waitgroup
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// channels for pipeline
	record := make(chan []string)
	done := make(chan bool)

	go transformer.ProcessRecord(wg, record, done)
	go p.Read(wg, record, done)

	wg.Wait()

	count := len(transformer.reporter.GetErrors())
	expected := 0
	if count != expected {
		t.Fatalf("Expected %d, got %d", expected, count)
	}

	count = transformer.reporter.GetTotalTransformedRecords()
	expected = 100
	if count != expected {
		t.Fatalf("Expected %d, got %d", expected, count)
	}
}

func TestProcessRecordFails(t *testing.T) {
	reporter := report.NewMockReporter(false)
	transformer := NewHTMLTransformer(reporter)

	path := utils.RootDir()
	p := parser.NewCSVParser(path+"/internal/testdata/fail_process_record.csv", reporter)

	// create waitgroup
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// channels for pipeline
	record := make(chan []string)
	done := make(chan bool)

	go transformer.ProcessRecord(wg, record, done)
	go p.Read(wg, record, done)

	wg.Wait()

	count := len(transformer.reporter.GetErrors())
	expected := 0
	if count == expected {
		t.Fatalf("Expected %d, got %d", expected, count)
	}

	count = transformer.reporter.GetTotalTransformedRecords()
	expected = 1
	if count < expected {
		t.Fatalf("Expected > %d, got %d", expected, count)
	}
}
