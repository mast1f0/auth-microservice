package seed

//seed for development
import (
	"auth-microservice/internal/adapters/crypto"
	"auth-microservice/internal/core/domain"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func RunSeed() {
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
		log.Fatal(err)
	}
	defer db.Close()

	bcryptHasher := crypto.NewBcryptHasher()
	users := []domain.User{
		{Login: "admin", Role: domain.RoleAdmin, HashedPwd: bcryptHasher.HashPassword("admin123")},
		{Login: "buyer", Role: domain.RoleBuyer, HashedPwd: bcryptHasher.HashPassword("buyer123")},
		{Login: "seller", Role: domain.RoleSeller, HashedPwd: bcryptHasher.HashPassword("seller123")},
	}
	for _, user := range users {
		result, err := db.Exec(
			`INSERT INTO users (login, role_id, password_hash, created_at)
			 SELECT $1, r.id, $2, $3
			 FROM roles r
			 WHERE r.name = $4
			 ON CONFLICT (login) DO NOTHING`,
			user.Login,
			user.HashedPwd,
			time.Now(),
			user.Role,
		)
		if err != nil {
			log.Println("error inserting user", err)
			continue
		}

		affected, err := result.RowsAffected()
		if err == nil && affected == 0 {
			log.Printf("user %s was skipped: login already exists or role is missing", user.Login)
		}
	}
}
