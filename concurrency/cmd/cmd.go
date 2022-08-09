package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

// set defaults
const (
	DefaultMaxReceivers int           = 1
	DefaultMaxRunTime   time.Duration = 2 * time.Second
)

var ErrorArgsDirSpecified = errors.New("A directory cannot be passed as an argument.")

// Args definition of all command line args
type Args struct {
	FilePath      string
	PrimaryTitle  string
	MaxRunTime    time.Duration
	NoOfReceivers int
}

func Process() {
	var args Args

	// accept args from stdin
	flag.StringVar(&args.FilePath, "filePath", "", "Absolute path to tsv file for processing.")
	flag.StringVar(&args.PrimaryTitle, "primaryTitle", "", "filter on primaryTitle")

	flag.IntVar(&args.NoOfReceivers, "noOfReceivers", DefaultMaxReceivers, "number of receivers to use in processing data")
	if args.NoOfReceivers <= 0 || args.NoOfReceivers > DefaultMaxReceivers {
		args.NoOfReceivers = DefaultMaxReceivers
	}

	flag.DurationVar(&args.MaxRunTime, "maxRunTime", DefaultMaxRunTime, "filter on maxRunTime")
	if args.MaxRunTime <= 0 {
		args.MaxRunTime = DefaultMaxRunTime
	}

	flag.Parse()

	// check if any arg is passed
	if len(os.Args[1:]) == 0 {
		flag.PrintDefaults()
		return
	}

	// get file's info
	f, err := os.Stat(args.FilePath)
	if err != nil {
		panic(err)
	}

	// check it is a file
	if f.IsDir() {
		panic(ErrorArgsDirSpecified)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(args.MaxRunTime))
	defer cancel()

	reader := NewReader(args.FilePath, args.NoOfReceivers)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go reader.Read(ctx, wg, args)

	wg.Wait()

	fmt.Println("Total Read: ", reader.totalProcessedRecords)
	fmt.Println("Total Failed: ", reader.totalFailedRecords)
	fmt.Println("Duration: ", fmt.Sprintf("%.2f", reader.duration)+"s")
}
