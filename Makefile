run:
	go run cmd/api/*.go

build:
	go build -o bin/application cmd/api/*.go
