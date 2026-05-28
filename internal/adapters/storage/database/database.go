package database

import (
	"auth-microservice/internal/core/domain"
	"auth-microservice/internal/core/ports"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	DB *sql.DB
}

func (db *Database) UserByID(id int64) (*domain.User, error) {
	var usr domain.User
	err := db.DB.QueryRow(
		`SELECT u.id, u.login, u.password_hash, r.name, u.created_at
		 FROM users u
		 JOIN roles r ON r.id = u.role_id
		 WHERE u.id = $1`,
		id,
	).Scan(&usr.Id, &usr.Login, &usr.HashedPwd, &usr.Role, &usr.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrUserNotFound
		}
		return nil, err
	}
	return &usr, nil
}

func (db *Database) UserExists(login string) bool {
	var exists bool
	err := db.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)`, login).Scan(&exists)
	if err != nil {
		log.Println("failed to check user existence")
		return false
	}
	return exists
}

func (db *Database) AddUser(user *domain.User) (*domain.User, error) {
	var (
		id        int64
		createdAt = time.Now()
	)

	err := db.DB.QueryRow(
		`INSERT INTO users (login, password_hash, role_id, created_at)
		 VALUES (
			$1,
			$2,
			(SELECT id FROM roles WHERE name = $3),
			$4
		 ) RETURNING id`,
		user.Login,
		user.HashedPwd,
		user.Role,
		createdAt,
	).Scan(&id)

	if err != nil {
		log.Println("failed to add user")
		return nil, err
	}
	return &domain.User{Id: id, Login: user.Login, Role: user.Role, CreatedAt: createdAt}, nil
}

func (db *Database) UserByLogin(login string) (*domain.User, error) {
	var usr domain.User
	err := db.DB.QueryRow(
		`SELECT u.id, u.login, u.password_hash, r.name, u.created_at
		 FROM users u
		 JOIN roles r ON r.id = u.role_id
		 WHERE u.login = $1`,
		login,
	).Scan(&usr.Id, &usr.Login, &usr.HashedPwd, &usr.Role, &usr.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrUserNotFound
		}
		return nil, err
	}

	return &usr, nil
}

func (db *Database) DeleteUser(id int64) error {
	_, err := db.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetAllUsers() ([]*domain.User, error) {
	users := make([]*domain.User, 0)
	query := `
        SELECT u.id, u.login, r.name, u.created_at
        FROM users u
        JOIN roles r ON u.role_id = r.id
        ORDER BY u.id
    `
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Println(err)
		return nil, ports.ErrFailedToLoad
	}
	defer rows.Close()
	for rows.Next() {
		var usr domain.User
		var role string
		if err := rows.Scan(&usr.Id, &usr.Login, &role, &usr.CreatedAt); err != nil {
			log.Println(err)
			return nil, ports.ErrFailedToLoad
		}
		usr.Role = domain.Role(role)
		users = append(users, &usr)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, ports.ErrFailedToLoad
	}
	return users, nil
}

func (db *Database) UpdateUser(userId int64, newRole domain.Role) error {
	query := `
        UPDATE users 
        SET role_id = (SELECT id FROM roles WHERE name = $1)
        WHERE id = $2
    `
	_, err := db.DB.Exec(query, newRole, userId)
	if err != nil {
		return ports.ErrFailedToUpdate
	}
	return nil
}

func NewDatabase() (*Database, error) {
	dbInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("pgx", dbInfo)

	if err != nil {
		return &Database{}, err
	}
	return &Database{
		DB: db,
	}, nil
}
