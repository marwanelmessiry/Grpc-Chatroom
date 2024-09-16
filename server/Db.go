package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender TEXT,
		recipient TEXT,
		content TEXT,
		timestamp DATETIME
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
		return nil, err
	}

	return db, nil
}
