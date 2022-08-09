package cmd

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"sync"
	"time"
)

type Reader struct {
	// channels for pipeline
	record chan []string // chan for record read
	done   chan bool     // chan to signal EOF for reader

	receiversCount int // number of receivers to spin up

	file                  string  // abs path to file
	totalProcessedRecords int     // total number of records processed
	totalFailedRecords    int     // total number of records failed
	errors                []error // errors found while processing
	duration              float64 // total duration of entire process
}

// NewReader creates a new reader for tsv file
//
// f is abs path to tsv file to be read.
// r is the number of imdb receivers to spun up in order to process the file.
// if r is 0 or -ve, default will be 1 receiver.
func NewReader(f string, r int) *Reader {
	// there has to be at least one receiver if val is invalid
	if r <= 0 {
		r = 1
	}

	return &Reader{
		file:           f,
		record:         make(chan []string),
		done:           make(chan bool),
		receiversCount: r,
	}
}

// ReadLevelTwo reads from the tsv file for level2 requirement
func (r *Reader) Read(ctx context.Context, wg *sync.WaitGroup, filters Args) {
	var (
		start = time.Now()
		flag  bool
	)

	cctx, cancel := context.WithCancel(ctx)

	defer func() {
		close(r.done)
		close(r.record)
		r.duration += float64(time.Since(start).Seconds())
		wg.Done()
		cancel()
	}()

	// open file for reading
	f, err := os.Open(r.file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// mini listener to listen for when a receiver has found a match
	// if so, this listener receives it on behalf of the reader and
	// let's the reader know so it halts reading by setting flag = true.
	go func() {
		<-r.done
		flag = true
	}()

	// spin up receivers
	wg.Add(r.receiversCount)
	for i := 1; i <= r.receiversCount; i++ {
		go NewReceiver().Receive(cctx, i, &LevelTwoReceiverParams{
			Filter: filters.PrimaryTitle,
			Done:   r.done,
			Record: r.record,
			Wg:     wg,
		})
	}

	// create csv reader
	reader := csv.NewReader(f)
	reader.Comma = '\t'

	// skip headers
	_, _ = reader.Read()

	// read rest of records
	for {
		select {
		// context cancelled/deadline exceeded
		case <-ctx.Done():
			flag = true
			return
		default:
			// read next record
			row, err := reader.Read()
			if err != nil {
				if err != io.EOF {
					r.errors = append(r.errors, err)
					r.totalFailedRecords++
					continue
				}

				flag = true
				break
			}

			// pass record to receivers
			if !flag {
				r.totalProcessedRecords++
				r.record <- row
				continue
			}
		}

		// either reading has completed, context cancelled/deadline exceeded, or
		// we have found a match by virtue of using a filter,
		// send a done signal to all receivers.
		if flag {
			r.done <- true
			break
		}
	}
}
