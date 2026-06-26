BIN := bin

.PHONY: build build-app build-db build-validate run test validate clean

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

## test: run all unit tests
test:
	go test ./...

## validate: build and run the accuracy validation against the database
validate: build-validate
	./$(BIN)/validate

## clean: remove compiled binaries
clean:
	rm -rf $(BIN)
