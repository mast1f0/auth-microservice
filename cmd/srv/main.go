package main

import (
	"auth-microservice/internal/adapters/http"
	"auth-microservice/internal/adapters/http/handlers"
	jwtutil "auth-microservice/internal/adapters/jwt"
	"auth-microservice/internal/adapters/storage/database"
	"auth-microservice/internal/core/service"
	seed2 "auth-microservice/internal/seed"
	"flag"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables")
	}
}
func main() {
	seed := flag.Bool("seed", false, "Seed the database")
	flag.Parse()

	if *seed {
		seed2.RunSeed()
		return
	}

	db, err := database.NewDatabase()
	if err != nil {
		log.Println("Не удалось подключиться к бд:", err)
		return
	}
	manager := jwtutil.Manager{
		Secret: []byte("superSecret"),
	}
	userService := service.NewService(db)
	handler := handlers.NewHandlers(userService, &manager)
	router := http.NewRouter(handler)
	srv := http.NewServer(router)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}
}
