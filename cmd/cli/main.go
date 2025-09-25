package main

import (
	"context"
	"log"
	"os"

	"zoubun/internal/term"
)

func main() {
	cmd := term.CLICommands()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
