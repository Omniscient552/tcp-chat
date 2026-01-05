package main

import (
	"fmt"
	"os"

	"tcp-chat/internal/models"
	"tcp-chat/internal/server"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("\x1b[33;1;7mUSAGE:\x1b[0;34;3;7m go run ./cmd/app/ \n\x1b[0;33;1;7mOR    \x1b[0;34;3;7m go run ./cmd/app/ \x1b[33m[PORT]\x1b[0m")
	}

	if len(os.Args) > 1 {
		models.PORT = ":" + os.Args[1]
	}
	server.RunServer()
}
