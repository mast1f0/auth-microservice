package database

import (
	"auth-microservice/internal/adapters/storage/migrations"
	"auth-microservice/internal/core/domain"
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

func (db *Database) UserByID(id uint) (*domain.User, error) {
	var usr domain.User
	_ = db.DB.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&usr.Id, &usr.Role, &usr.Login, &usr.HashedPwd, &usr.CreatedAt)
	if usr.Id == 0 {
		return &domain.User{}, errors.New("Нет такого пользователя")
	}
	return &usr, nil
}

func (db *Database) UserExists(login string) bool {
	query, err := db.DB.Query(`SELECT * FROM users`)
	if err != nil {
		log.Println("не удалось сделать выборку")
	}

	var usr domain.User

	for query.Next() {
		if query.Scan(&usr.Id, &usr.Login, &usr.HashedPwd, &usr.CreatedAt); usr.Login == login {
			return true
		}
	}
	return false
}

func (db *Database) AddUser(user *domain.User) (*domain.User, error) {
	var id int64
	err := db.DB.QueryRow(
		`INSERT INTO users (login, password_hash, role, created_at) 
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Login,
		user.HashedPwd,
		user.Role,
		time.Now(),
	).Scan(&id)

	if err != nil {
		log.Println("Не удалось добавить пользователя")
		return nil, err
	}
	return &domain.User{Id: id, Login: user.Login, Role: user.Role}, nil
}

func (db *Database) UserByLogin(login string) (*domain.User, error) {
	var usr domain.User
	err := db.DB.QueryRow(
		`SELECT id, login, password_hash, role, created_at FROM users WHERE login = $1`,
		login,
	).Scan(&usr.Id, &usr.Login, &usr.HashedPwd, &usr.Role, &usr.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Пользователь не найден")
		}
		return nil, err
	}

	return &usr, nil
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
	fmt.Println(db_info)
	db, err := sql.Open("pgx", db_info)

	if err != nil {
		return &Database{}, err
	}
	err = migrations.RunMigrations(db)
	if err != nil {
		log.Fatal(err)
		return &Database{}, err
	}
	return &Database{
		DB: db,
	}, nil
}
