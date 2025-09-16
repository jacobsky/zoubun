BINARY_SERVER=bin\server.exe
BINARY_CLI=bin\cli.exe

build:
	sqlc generate
	go build -o ${BINARY_SERVER} cmd/server/main.go
	go build -o ${BINARY_CLI} cmd/cli/main.go

serve:
	docker-compose up --build
