#Makefile
include .env

up:
	docker-compose up -d

up-db:
	docker-compose up -d db

up-redis:
	docker-compose up -d redis

down:
	docker-compose down

logs:
	docker-compose logs -f

psql:
	docker exec -it $(NAME_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

migrate:
	go run cmd/api/main.go migrate

reset:
	docker-compose down -v
	docker-compose up -d

env:
	@echo "DB_USER=$(DB_USER)"
	@echo "DB_PASSWORD=$(DB_PASSWORD)"
	@echo "DB_NAME=$(DB_NAME)"