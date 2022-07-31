package parser

import (
	"sync"
	"testing"

	"github.com/dele454/medium/csv-transform-to-html/internal/report"
	"github.com/dele454/medium/csv-transform-to-html/internal/utils"
)

func TestCSVRead(t *testing.T) {
	reporter := report.NewMockReporter()

	path := utils.RootDir()
	p := NewCSVParser(path+"/internal/testdata/100_sales_records.csv", reporter)

	// create waitgroup
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// channels for pipeline
	record := make(chan []string)
	done := make(chan bool)
	expected := 100

	var counter int
	var flag bool

	go func() {
		defer wg.Done()

		for {
			select {
			case <-record:
				counter++
			case <-done:
				flag = true
			}

			if flag {
				break
			}
		}
	}()

	go p.Read(wg, record, done)
	wg.Wait()

	if counter != expected {
		t.Fatalf("\nFormat Mismatch:\nExpected: %v\nGot: %v", expected, counter)
	}
}
