package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

const copySufix = "_export"

type storage struct {
	cnn       *sql.DB
	tableName string
	query     string
}

func newStorage(query string) (*storage, error) {
	cnn, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SCHEMA"),
	))

	if err != nil {
		return nil, err
	}

	tableName := strings.ReplaceAll(uuid.New().String()+copySufix, "-", "")

	db := &storage{
		cnn:       cnn,
		query:     query,
		tableName: tableName,
	}

	return db, nil
}

func (s *storage) copyData() error {
	createTable := fmt.Sprintf(
		"CREATE TABLE %s AS %s",
		s.tableName,
		s.query,
	)

	if _, err := s.cnn.Exec(createTable); err != nil {
		return err
	}

	return nil
}

func (s *storage) dropTable() error {
	dropTable := fmt.Sprintf(
		"DROP TABLE %s",
		s.tableName,
	)

	_, err := s.cnn.Exec(dropTable)

	return err
}
