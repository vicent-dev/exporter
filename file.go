package main

import (
	"bufio"
	"fmt"
	"github.com/en-vee/alog"
	"os"
	"runtime"
	"sync"
	"time"
)

func createFile(directory string, filename string) (*os.File, error, string) {
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, err, ""
	}

	path := fmt.Sprintf("%s/%s_%d.csv",
		directory,
		filename,
		time.Now().Unix(),
	)
	f, err := os.Create(path)

	return f, err, path
}

func writeFileByChunks(f *os.File, s *storage) error {
	wg := sync.WaitGroup{}
	maxWg := runtime.GOMAXPROCS(0)
	c, err := s.getCount()

	if err != nil {
		return err
	}

	linesPerChunk := c / maxWg

	// for small exports we don't need more than one goroutines
	if c < 100 {
		linesPerChunk = c
	}

	// we want to force the first chunk processed first to add headers
	wg.Add(1)
	writeChunk(f, 0, linesPerChunk, s, &wg)

	if linesPerChunk != c {
		for i := 1; i <= maxWg; i++ {
			wg.Add(1)
			go writeChunk(f, i*linesPerChunk, linesPerChunk, s, &wg)
		}
	}

	wg.Wait()

	alog.Info("Count records exported: %d", c)

	return nil
}

func writeChunk(f *os.File, start, size int, s *storage, wg *sync.WaitGroup) {
	defer wg.Done()
	w := bufio.NewWriter(f)

	if err := s.writeData(w, start, size); err != nil {
		alog.Error(err.Error())
		return
	}

	if err := w.Flush(); err != nil {
		alog.Error(err.Error())
		return
	}
}
