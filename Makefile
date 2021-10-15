run:
	go run cmd/api/*.go

build:
	go build -o bin/server cmd/api/*.go

build-linux:
	env GOOS=linux GOARCH=amd64 go build -o bin/server cmd/api/*.go
	
archive:
	git archive -o bin/archived.zip HEAD