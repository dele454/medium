package cmd

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

// using sync.Pool to manage Movie structs
var pool = sync.Pool{
	New: func() any {
		return new(Movie)
	},
}

// IMDbReceiver handles the processing of movie titles
type IMDbReceiver struct{}

// LevelTwoReceiverParams for LevelTwoReceiver func
// As a way of compacting the params for LevelTwoReceiver
type LevelTwoReceiverParams struct {
	Wg     *sync.WaitGroup // waitgroup
	Record <-chan []string // record read from file
	Done   chan bool       // signal reading has completed or match found (if filter is applied)
	Filter string          // single filter to apply to record
}

// NewReceiver gets a new receiver
func NewReceiver() *IMDbReceiver {
	return new(IMDbReceiver)
}

// LevelTwoReceiver level two receiver
func (r *IMDbReceiver) Receive(ctx context.Context, id int, params *LevelTwoReceiverParams) {
	var counter int

	defer func() {
		params.Wg.Done()
		fmt.Println("Receiver #", id, fmt.Sprintf(" Received %v titles", counter))
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case rec := <-params.Record:
			// if empty ignore
			if len(rec) == 0 {
				continue
			}

			counter++

			// unmarshal record
			movie := NewMovie()
			movie.ReturnMovieToPool()
			movie = r.Unmarshal(rec, movie)

			// if filter applied doesn't match, move on
			if params.Filter != "" && movie.PrimaryTitle != params.Filter {
				continue
			}

			// user didnt supply a filter, move on
			if params.Filter == "" {
				continue
			}

			fmt.Println("Receiver #", id, fmt.Sprintf(" Found %+v", movie))
			params.Done <- true
			break
		case <-params.Done:
			return
		}
	}
}

// Movie basic details of a movie
type Movie struct {
	Tconst         string `tsv:"tconst"`
	TitleType      string `tsv:"titleType" processor:"required"`
	PrimaryTitle   string `tsv:"primaryTitle" processor:"required"`
	OriginalTitle  string `tsv:"originalTitle" processor:"required"`
	IsAdult        string `tsv:"isAdult" processor:"required"`
	StartYear      string `tsv:"startYear" processor:"year,required"`
	EndYear        string `tsv:"endYear"`
	RuntimeMinutes string `tsv:"runtimeMinutes" processor:"required"`
	Genres         string `tsv:"genres" processor:"required"`
}

// NewMovie returns a new Movie struct from a pool
func NewMovie() Movie {
	return *pool.Get().(*Movie)
}

// ReturnMovieToPool returns a movie struct back into the pool
func (m *Movie) ReturnMovieToPool() {
	pool.Put(m)
}

// Unmarshal unmarshals records found into the Movie struct
func (r *IMDbReceiver) Unmarshal(record []string, movie Movie) Movie {
	s := reflect.ValueOf(movie).Type()
	for i := 0; i < s.NumField(); i++ {
		reflect.ValueOf(&movie).Elem().Field(i).SetString(record[i])
	}

	return movie
}
