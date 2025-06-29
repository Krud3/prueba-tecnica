#Makefile
include .env

#up all docker services
up:
	docker-compose up -d

#only db
up-db:
	docker-compose up -d db

#only redis
up-redis:
	docker-compose up -d redis

#shut-down services
down:
	docker-compose down

#logs
logs:
	docker-compose logs -f

#enter psql tool
psql:
	docker exec -it $(NAME_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

#migrate manually	
migrate:
	migrate -path migrations/ -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

#fix migrate due to a bad migrate
migrate-fix:
	migrate -path migrations/ -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" force 1

#down migrate
migrate-down:
	migrate -path migrations/ -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

#run app
run-app:
	go run cmd/api/main.go

#reset docker services
reset:
	docker-compose down -v
	docker-compose up -d

#get .evn values
env:
	@echo "DB_USER=$(DB_USER)"
	@echo "DB_PASSWORD=$(DB_PASSWORD)"
	@echo "DB_NAME=$(DB_NAME)"