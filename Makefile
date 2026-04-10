.PHONY: run-dev run-seed migrate-up migrate-down

DB_URL ?= postgres://auth_user:strong_password@localhost:5432/auth_db?sslmode=disable

run-dev:
	go run cmd/srv/main.go

run-seed:
	go run cmd/srv/main.go --seed

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1
