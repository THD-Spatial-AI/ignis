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

## create-db: create the database if absent (or drop only the tabula schema if it exists), then populate via build_db
create-db: build-db
	@source <(sed 's/[[:space:]]*#.*//' .env | grep -v '^[[:space:]]*$$') && \
	if PGPASSWORD=$$DB_PASSWORD psql -h $$DB_HOST -p $$DB_PORT -U $$DB_USER -d postgres \
	    -tc "SELECT 1 FROM pg_database WHERE datname='$$DB_NAME'" | grep -q 1; then \
	    echo "Database '$$DB_NAME' exists — dropping tabula schema only." && \
	    PGPASSWORD=$$DB_PASSWORD psql -h $$DB_HOST -p $$DB_PORT -U $$DB_USER -d $$DB_NAME \
	        -c "DROP SCHEMA IF EXISTS tabula CASCADE;"; \
	else \
	    PGPASSWORD=$$DB_PASSWORD psql -h $$DB_HOST -p $$DB_PORT -U $$DB_USER -d postgres \
	        -c "CREATE DATABASE $$DB_NAME;" && echo "Created database '$$DB_NAME'."; \
	fi && \
	./$(BIN)/build_db

## test: run all unit tests
test:
	go test ./...

## validate: build and run the accuracy validation against the database
validate: build-validate
	./$(BIN)/validate

## clean: remove compiled binaries
clean:
	rm -rf $(BIN)
