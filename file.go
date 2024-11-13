package main

import (
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

func writeFileByChunks(file *os.File, database *storage) error {
	wg := sync.WaitGroup{}

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		wg.Done()
	}

	wg.Wait()
	return nil
}
