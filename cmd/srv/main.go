package main

import (
	"auth-microservice/internal/adapters/crypto"
	server "auth-microservice/internal/adapters/http"
	"auth-microservice/internal/adapters/http/handlers"
	jwtutil "auth-microservice/internal/adapters/jwt"
	"auth-microservice/internal/adapters/storage/database"
	"auth-microservice/internal/core/service"
	seed2 "auth-microservice/internal/seed"
	"os"

	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

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
		Secret: []byte(os.Getenv("JWT_SECRET")),
	}
	bcryptHasher := crypto.NewBcryptHasher()
	userService := service.NewService(db, bcryptHasher)
	handler := handlers.NewHandlers(userService, &manager)
	router := server.NewRouter(handler)
	go func() { log.Println(http.ListenAndServe("localhost:6060", nil)) }()
	srv := server.NewServer(router)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}
}
