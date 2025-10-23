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

		createTableSQL := `
		CREATE TABLE IF NOT EXISTS public.listings (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			listing_type VARCHAR(10) NOT NULL CHECK (listing_type IN ('rent', 'sale')),
			price INT NOT NULL CHECK (price > 0),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);
		`

		_, err = dbSql.Exec(createTableSQL)
		if err != nil {
			panic(fmt.Sprintf("failed to create listings table: %v", err))
		}
	})

	return dbSql
}
