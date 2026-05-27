.PHONY: run-dev run-seed migrate-up migrate-down pprof test

DB_URL ?= postgres://postgres:postgres123@localhost:5432/auth_db?sslmode=disable

run-dev:
	go run cmd/srv/main.go

run-seed:
	go run cmd/srv/main.go --seed

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

pprof:
	go tool pprof -http=:6061 http://localhost:6060/debug/pprof/profile

test:
	go test ./...
