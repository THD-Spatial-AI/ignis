SHELL := /bin/bash
BIN   := bin

.PHONY: build build-app build-db build-validate run create-db test validate clean

## build: compile all binaries into bin/
build: build-app build-db build-validate

## build-app: compile the HTTP API server
build-app:
	go build -o $(BIN)/app cmd/app/main.go

## build-db: compile the database loader (destructive — drops and recreates tables)
build-db:
	go build -o $(BIN)/build_db cmd/build_db/main.go

## build-validate: compile the accuracy validation tool
build-validate:
	go build -o $(BIN)/validate cmd/validate/main.go

## run: build and start the HTTP API server
run: build-app
	./$(BIN)/app

## create-db: create the PostgreSQL database named in .env (DB_NAME); run once before build-db
create-db:
	@source <(sed 's/[[:space:]]*#.*//' .env | grep -v '^[[:space:]]*$$') && \
	PGPASSWORD=$$DB_PASSWORD psql -h $$DB_HOST -p $$DB_PORT -U $$DB_USER -d postgres \
	    -c "CREATE DATABASE $$DB_NAME;" && \
	echo "Database '$$DB_NAME' created."

## test: run all unit tests
test:
	go test ./...

## validate: build and run the accuracy validation against the database
validate: build-validate
	./$(BIN)/validate

## clean: remove compiled binaries
clean:
	rm -rf $(BIN)
