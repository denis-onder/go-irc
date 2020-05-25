package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func connect() (*sql.DB, error) {
	connStr := "postgres://postgres:postgres@localhost:5432/goirc"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}
