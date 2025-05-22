DB_URL="postgres://admin:admin@localhost:5436/postgres?sslmode=disable"

migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down

migrate-new:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $${name}


.PHONY: up down restart logs

up:
	docker-compose up -d

down:
	docker-compose down -v

restart: down up

logs:
	docker-compose logs -f --tail=100
