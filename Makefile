BINARY_SERVER=bin\server.exe
BINARY_CLI=bin\cli.exe

build:
	sqlc generate
	go build -o ${BINARY_SERVER} cmd/server/main.go
	go build -o ${BINARY_CLI} cmd/cli/main.go

serve:
	tern migrate
	go run ./cmd/server/

docker:
	docker build . --platform=linus/amd64
	docker build . --platform=linux/arm64
