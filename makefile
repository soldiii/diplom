.PHONY: build
build:
	go build -v ./cmd/main.go

migrate-up: 
	migrate -path migrations -database "postgres://postgres:postgres@localhost/supervisor_app_bd?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost/supervisor_app_bd?sslmode=disable" down


