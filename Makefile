#Makefile
include .env

up:
	docker-compose up -db

down:
	docker-compose down

logs:
	docker-compose logs -f

psql:
	docker exec -it $(CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME)

migrate:
	go run cmd/api/main.go migrate

reset:
	docker-compose down -v
	docker-compose up -d

env:
	@echo "DB_USER=$(DB_USER)"
	@echo "DB_PASSWORD=$(DB_PASSWORD)"
	@echo "DB_NAME=$(DB_NAME)"