.PHONY: build
build:
	docker-compose build app

run:
	docker-compose up app

migrate-up: 
	docker-compose run -v C:\Users\malax\OneDrive\Рабочий\стол\diplom\migrations:/migrations --network host migrate/migrate -path ./migrations -database "postgresql://postgres:postgres@database/supervisor_app_bd?sslmode=disable" 
migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@database/supervisor_app_bd?sslmode=disable" down


