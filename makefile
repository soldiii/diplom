.PHONY: build
build:
	go build -v ./cmd/main.go

build-docker:
	docker-compose build

run:
	./main.exe

run-docker:
	docker-compose up

migrate-up: 
	migrate -path migrations -database "postgres://postgres:postgres@diplom-database-1:5432/supervisor_app_bd?sslmode=disable" up 
migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@diplom-database-1:5432/supervisor_app_bd?sslmode=disable" down


