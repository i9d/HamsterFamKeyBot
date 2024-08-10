package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func initializeDB() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	createTables := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        chat_id BIGINT UNIQUE NOT NULL,
        user_id BIGINT NOT NULL,
        username TEXT,
        firstname TEXT,
        last_check TIMESTAMP,
        user_limit INT DEFAULT 1,
        already_got INT DEFAULT 0
    );
    CREATE TABLE IF NOT EXISTS codes (
        id SERIAL PRIMARY KEY,
        code TEXT NOT NULL,
        used BOOLEAN DEFAULT FALSE,
        used_by BIGINT,
        game_type TEXT
    );`

	_, err = db.Exec(createTables)
	if err != nil {
		return nil, err
	}

	return db, nil
}
