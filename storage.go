package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/en-vee/alog"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

	tableName := strconv.Itoa(int(time.Now().Unix())) + copySufix

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

func (s *storage) getCount() (int, error) {
	count := fmt.Sprintf("SELECT COUNT(*) FROM %s", s.tableName)

	rows, err := s.cnn.Query(count)

	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var c int

	rows.Next()

	if err = rows.Scan(&c); err != nil {
		return 0, nil
	}

	return c, nil
}

func (s *storage) writeData(w *bufio.Writer, start, size int) error {
	data := fmt.Sprintf("SELECT * FROM %s LIMIT %d,%d", s.tableName, start, size)
	alog.Info("Exporting records from %d to %d.", start, start+size)
	rows, err := s.cnn.Query(data)

	if err != nil {
		return err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			alog.Error(err.Error())
		}
	}(rows)

	cols, _ := rows.Columns()
	vals := make([]interface{}, len(cols))
	for i, _ := range cols {
		vals[i] = new(sql.RawBytes)
	}

	if start == 0 {
		if _, err := w.WriteString(strings.Join(cols, ";") + "\n"); err != nil {
			return err
		}
	}

	for rows.Next() {
		if err := rows.Scan(vals...); err != nil {
			return err
		}

		var lsb strings.Builder
		for i, v := range vals {
			lsb.WriteString(string(*(v.(*sql.RawBytes))))
			if i != len(vals)-1 {
				lsb.WriteString(";")
			}
		}

		lsb.WriteString("\n")
		if _, err := w.WriteString(lsb.String()); err != nil {
			return err
		}
	}

	return nil
}
