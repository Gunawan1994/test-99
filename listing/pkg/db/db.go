package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	dbSql  *sql.DB
	oncePg sync.Once
)

func NewConn() *sql.DB {
	var err error

	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"postgres",
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_ADDR"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DATABASE"),
	)

	oncePg.Do(func() {
		dbSql, err = sql.Open("postgres", conn)
		if err != nil {
			panic(err)
		}

		if err := dbSql.Ping(); err != nil {
			panic(fmt.Sprintf("failed to connect to PostgreSQL: %v", err))
		}
	})

	return dbSql
}
