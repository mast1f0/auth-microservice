package migrations

import (
	"database/sql"
)

func RunMigrations(db *sql.DB) error {
	db.Exec(`DROP TABLE IF EXISTS users;`)
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users  (
    id SERIAL PRIMARY KEY,
    role VARCHAR(20) NOT NULL,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);`)
	return err
}
