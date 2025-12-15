# Standart MakeFile

# variables
APP_NAME=tradelog
DSN="host=localhost user=admin password=secret dbname=tradelog port=5432 sslmode=disable"

run:
	go run cmd/api/main.go

build:
	go build -o bin/$(APP_NAME) cmd/api/main.go

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

clean:
	rm -f bin/$(APP_NAME)
	rm -rf tmp

dev:
	air