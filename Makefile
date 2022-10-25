build:
	docker-compose build

run:
	docker-compose up

migrate-up:
	migrate -path ./schema -database 'postgres://admin:admin@localhost:5400/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./schema -database 'postgres://admin:admin@localhost:5400/postgres?sslmode=disable' down