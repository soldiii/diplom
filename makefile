.PHONY: build
build:
	go build -v ./cmd/main.go

build-docker:
	docker-compose build app

run:
	./main.exe

run-docker:
	docker-compose up app

migrate-up: 
	migrate -path migrations -database "postgres://postgres:postgres@localhost/supervisor_app_bd?sslmode=disable" up 
migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost/supervisor_app_bd?sslmode=disable" down


