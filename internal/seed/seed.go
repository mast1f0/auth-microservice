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

	users := []domain.User{
		{Login: "admin", Role: domain.RoleAdmin, HashedPwd: crypto.HashPassword("admin123")},
		{Login: "buyer", Role: domain.RoleBuyer, HashedPwd: crypto.HashPassword("buyer123")},
		{Login: "seller", Role: domain.RoleSeller, HashedPwd: crypto.HashPassword("seller123")},
	}
	for _, user := range users {
		_, err := db.Exec("INSERT INTO users (login, role, password_hash, created_at) VALUES ($1, $2, $3, $4)  ON CONFLICT (login) DO NOTHING", user.Login, user.Role, user.HashedPwd, time.Now())
		if err != nil {
			log.Println("error inserting user", err)
		}
	}
}
