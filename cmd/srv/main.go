package main

import (
	"auth-microservice/internal/database"
	"auth-microservice/internal/server"
	"log"
)

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Println("Не удалось подключиться к бд")
	}

	srv := server.NewServer(db)
	srv.RunServer()
}
