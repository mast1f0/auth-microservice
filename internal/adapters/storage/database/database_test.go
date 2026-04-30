package database

import (
	"auth-microservice/internal/core/domain"
	"database/sql"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func newMockDB(t *testing.T) (*Database, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New failed: %v", err)
	}

	cleanup := func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sql expectations: %v", err)
		}
		db.Close()
	}

	return &Database{DB: db}, mock, cleanup
}

func TestUserByID(t *testing.T) {
	repo, mock, cleanup := newMockDB(t)
	defer cleanup()

	createdAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "login", "password_hash", "name", "created_at"}).
		AddRow(int64(1), "user", []byte("hash"), domain.RoleBuyer, createdAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, u.login, u.password_hash, r.name, u.created_at
		 FROM users u
		 JOIN roles r ON r.id = u.role_id
		 WHERE u.id = $1`)).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	user, err := repo.UserByID(1)
	if err != nil {
		t.Fatalf("UserByID returned error: %v", err)
	}
	if user.Id != 1 || user.Login != "user" || user.Role != domain.RoleBuyer {
		t.Fatalf("unexpected user returned: %+v", user)
	}
}

func TestUserByIDNotFound(t *testing.T) {
	repo, mock, cleanup := newMockDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, u.login, u.password_hash, r.name, u.created_at
		 FROM users u
		 JOIN roles r ON r.id = u.role_id
		 WHERE u.id = $1`)).
		WithArgs(int64(1)).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.UserByID(1)
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
	if err != domain.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestAddUser(t *testing.T) {
	repo, mock, cleanup := newMockDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (login, password_hash, role_id, created_at)
		 VALUES (
			$1,
			$2,
			(SELECT id FROM roles WHERE name = $3),
			$4
		 ) RETURNING id`)).
		WithArgs("user", []byte("hash"), domain.RoleAdmin, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(10)))

	user, err := repo.AddUser(&domain.User{Login: "user", HashedPwd: []byte("hash"), Role: domain.RoleAdmin})
	if err != nil {
		t.Fatalf("AddUser returned error: %v", err)
	}
	if user.Id != 10 || user.Login != "user" || user.Role != domain.RoleAdmin {
		t.Fatalf("unexpected user returned: %+v", user)
	}
	if user.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}
}

func TestUserByLoginNotFound(t *testing.T) {
	repo, mock, cleanup := newMockDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, u.login, u.password_hash, r.name, u.created_at
		 FROM users u
		 JOIN roles r ON r.id = u.role_id
		 WHERE u.login = $1`)).
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.UserByLogin("missing")
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
	if err != domain.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	repo, mock, cleanup := newMockDB(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = $1")).WithArgs(int64(5)).WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.DeleteUser(5); err != nil {
		t.Fatalf("DeleteUser returned error: %v", err)
	}
}

func TestUserExists(t *testing.T) {
	repo, mock, cleanup := newMockDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)`)).
		WithArgs("user").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	if !repo.UserExists("user") {
		t.Fatal("expected user to exist")
	}
}
