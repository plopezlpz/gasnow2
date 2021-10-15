run:
	go run cmd/api/*.go

build:
	go build -o bin/server cmd/api/*.go

build-linux:
	env GOOS=linux GOARCH=amd64 go build -o bin/application cmd/api/*.go

archive:
	zip -r app.zip bin -x bin/.DS_Store