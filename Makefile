BINARY_SERVER=bin\server.exe
BINARY_CLI=bin\cli.exe

build:
	go build -o ${BINARY_SERVER} cmd/server/main.go
	go build -o ${BINARY_CLI} cmd/cli/main.go

serve:
	build
	go run ${BINARY_SERVER}

