package main

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"runtime"
	"sync"
)

func createFile(directory string) (*os.File, error) {
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, err
	}

	return os.Create(directory + "/" + uuid.New().String() + ".csv")
}

func writeFileByChunks(f *os.File, s *storage) error {
	wg := sync.WaitGroup{}
	maxWg := runtime.GOMAXPROCS(0)
	c, err := s.getCount()

	if err != nil {
		return err
	}

	linesPerWg := c / maxWg
	extraLines := c % maxWg

	if extraLines > 0 {
		linesPerWg = c / (maxWg - 1)
	}

	fmt.Println(linesPerWg)
	for i := 0; i < maxWg; i++ {
		wg.Add(1)

		wg.Done()
	}

	wg.Wait()
	return nil
}
