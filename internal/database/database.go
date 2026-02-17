package database

import (
	"auth-microservice/internal/crypto"
	"auth-microservice/internal/models"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib" // Import the stdlib adapter
)

type Database struct {
	DB *sql.DB
}

func (db *Database) AddUser(user models.User) {
	user.HashedPwd = crypto.HashPassword(user.HashedPwd)
	db.DB.Exec(``, user.Id, user.Login, user.HashedPwd)
}

func (db *Database) UserByLogin(login string) *models.User {
	return &models.User{}
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("pgx", "user=postgres password=postgres123 dbname=... sslmode=disable")

	if err != nil {
		return &Database{}, err
	}

	return &Database{
		DB: db,
	}, nil
}
