package migrations

import (
	"database/sql"
)

func RunMigrations(db *sql.DB, path string) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users  (
    id SERIAL PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);`)
	return err
}
