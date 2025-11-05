package main

import (
	"errors"
	"flag"
	"runtime"
	"sync"

	"github.com/en-vee/alog"
)

func run() error {
	// process flags and load env vars
	loadEnv()

	q := flag.String("q", "", "Query to fetch data to be exported")
	directory := flag.String("d", ".", "Directory to save the exported file")
	filename := flag.String("f", "data", "Filename exported file")

	flag.Parse()

	if *q == "" {
		return errors.New("Query parameter [q] is mandatory")
	}

	s, err := newStorage(*q)

	if err != nil {
		return err
	}

	// take a snapshot of the data to be exported so it doesn't change from the beginning to the end of the process
	if err := s.copyData(); err != nil {
		return err
	} else {
		alog.Info("Snapshot created.")
	}

	defer func(s *storage) {
		if err := s.dropTable(); err != nil {
			alog.Error(err.Error())
		} else {
			alog.Info("Snapshot removed.")
		}
	}(s)

	nThreads := runtime.GOMAXPROCS(0) - 1
	fm := newFilesManager(*directory, *filename, nThreads)

	if err != nil {
		return err
	}

	if err := processInBatches(s, fm); err != nil {
		return err
	} else {
		alog.Info("Data was exported: %s", fm.mainFilePath)
	}

	return nil
}

func processInBatches(s *storage, fm *filesManager) error {
	wg := sync.WaitGroup{}

	c64, err := s.getCount()

	c := int(c64)

	if c == 0 {
		alog.Warn("No rows found based on criteria")
		return nil
	}

	if err != nil {
		return err
	}

	linesPerChunk := c / fm.nThreads

	// for small exports we don't need more than one goroutines
	if c < 100 {
		linesPerChunk = c
		wg.Add(1)
		go s.extractChunk(linesPerChunk, &wg, 0, fm)
	} else {
		for i := 0; i <= fm.nThreads; i++ {
			wg.Add(1)
			go s.extractChunk(linesPerChunk, &wg, i, fm)
		}
	}

	wg.Wait()

	totalRows, _ := fm.mergePartFiles()
	// fm.removePartFiles()

	alog.Info("Count records exported: %d", totalRows)

	return nil
}

func main() {
	if err := run(); err != nil {
		alog.Error(err.Error())
	}
}
