package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Init opens (or creates) the SQLite database and ensures tables exist.
func Init() {
	// Use DB_PATH environment variable if set, otherwise default to ./findit.db
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./findit.db"
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	migrate()
	log.Println("database ready")
}

func migrate() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id       TEXT PRIMARY KEY,
		email    TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`

	itemTable := `
	CREATE TABLE IF NOT EXISTS items (
		id          TEXT PRIMARY KEY,
		user_id     TEXT NOT NULL,
		type        TEXT NOT NULL CHECK(type IN ('lost','found')),
		name        TEXT NOT NULL,
		description TEXT,
		location    TEXT,
		date        TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	if _, err := DB.Exec(userTable); err != nil {
		log.Fatalf("failed to create users table: %v", err)
	}
	if _, err := DB.Exec(itemTable); err != nil {
		log.Fatalf("failed to create items table: %v", err)
	}
}
