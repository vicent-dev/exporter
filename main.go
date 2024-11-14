package main

import (
	"errors"
	"flag"
	"github.com/en-vee/alog"
)

func run() error {
	// process flags and load env vars
	err := loadEnv()
	if err != nil {
		return err
	}

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

	f, err, path := createFile(*directory, *filename)
	if err != nil {
		return err
	}

	if err := writeFileByChunks(f, s); err != nil {
		return err
	} else {
		alog.Info("Data was exported: %s", path)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		alog.Error(err.Error())
	}
}
