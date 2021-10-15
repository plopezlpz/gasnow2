run:
	go run cmd/api/*.go

build:
	go build -o bin/server cmd/api/*.go