.PHONY: run-dev run-seed

run-dev:
	go run cmd/srv/main.go

run-seed:
	go run cmd/srv/main.go --seed