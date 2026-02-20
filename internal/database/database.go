package database

import (
	"auth-microservice/internal/crypto"
	"auth-microservice/internal/migrations"
	"auth-microservice/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Import the stdlib adapter
)

type Database struct {
	DB *sql.DB
}

func (db *Database) UserCount() uint {
	var count uint = 0
	_ = db.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	return count
}

func (db *Database) UserByID(id uint) (*models.User, error) {
	var usr *models.User
	_ = db.DB.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&usr.Id, &usr.Login, &usr.HashedPwd, &usr.CreatedAt)
	if usr.Id == 0 {
		return &models.User{}, errors.New("Нет такого пользователя")
	}
	return usr, nil
}

func (db *Database) UserExists(login string) bool {
	query, err := db.DB.Query(`SELECT * FROM users`)
	if err != nil {
		log.Println("не удалось сделать выборку")
	}

	var usr models.User

	for query.Next() {
		if query.Scan(&usr.Id, &usr.Login, &usr.HashedPwd, &usr.CreatedAt); usr.Login == login {
			return true
		}
	}
	return false
}

func (db *Database) AddUser(user *models.User) error {
	HashedPwd := crypto.HashPassword(user.HashedPwd)
	_, err := db.DB.Exec(`INSERT INTO users (login, password_hash, created_at) VALUES ($1, $2, $3)`, user.Login, HashedPwd, time.Now())
	if err != nil || db.UserExists(user.Login) {
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
	db_info := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("pgx", db_info)

	if err != nil {
		return &Database{}, err
	}
	migrations.RunMigrations(db, "/internal/sql/create.sql")
	return &Database{
		DB: db,
	}, nil
}
