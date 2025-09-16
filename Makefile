BINARY_SERVER=bin/server.exe
BINARY_CLI=bin/cli.exe

build:
	sqlc generate
	go build -o ${BINARY_SERVER} cmd/server/main.go
	go build -o ${BINARY_CLI} cmd/cli/main.go

dev:
	docker-compose -f docker-compose.yaml -f docker-compose-dev.yaml up --build

prod:
	docker-compose up --build
