package database

import (
	"auth-microservice/internal/crypto"
	"auth-microservice/internal/migrations"
	"auth-microservice/internal/models"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Import the stdlib adapter
)

type Database struct {
	DB *sql.DB
}

// сделать проверку что пользователя нет
func (db *Database) UserExists() bool {

	return false
}

func (db *Database) AddUser(user *models.User) error {
	HashedPwd := crypto.HashPassword(user.HashedPwd)
	_, err := db.DB.Exec(`INSERT INTO users (login, password_hash, created_at) VALUES ($1, $2, $3)`, user.Login, HashedPwd, time.Now())
	if err != nil {
		log.Println("Не удалось добавить пользователя")
		return err
	}
	return nil
}

func (db *Database) UserByLogin(login string) (*models.User, error) {
	rows, err := db.DB.Query(`SELECT * FROM users`)
	if err != nil {
		log.Println(err)
	}
	var usr models.User
	for rows.Next() {
		if err = rows.Scan(&usr.Id, &usr.Login, &usr.HashedPwd, &usr.CreatedAt); usr.Login == login {
			return &usr, nil
		}
	}
	return &models.User{}, errors.New("Пользователь не найден!")
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("pgx", "user=auth_user password=strong_password dbname=auth_db sslmode=disable")

	if err != nil {
		return &Database{}, err
	}
	migrations.RunMigrations(db, "../../internal/sql/create.sql")
	return &Database{
		DB: db,
	}, nil
}
