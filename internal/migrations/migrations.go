package migrations

import (
	"database/sql"
	"os"
)

func RunMigrations(db *sql.DB, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(data))
	return err
}
