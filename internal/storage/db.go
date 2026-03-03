package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func InitDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {

		return nil, fmt.Errorf("failed to reach database: %w", err)
	}

	//Limits how many connections are oprn at once. Too many open connections can overwhelm the database
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	fmt.Println("Database connected successfully")
	return db, nil
}