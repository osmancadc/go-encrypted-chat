package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/osmancadc/go-encrypted-chat/internal/model"
	"github.com/osmancadc/go-encrypted-chat/internal/websocket"
	"github.com/osmancadc/go-encrypted-chat/pkg/logger"
)

var log = logger.NewLogger("INFO")

func main() {
	serverMode := flag.Bool("server", false, "Run in server mode")
	clientMode := flag.Bool("client", false, "Run in client mode")
	username := flag.String("user", "", "Username for client")
	flag.Parse()

	if *serverMode && *clientMode {
		fmt.Println("Error: Cannot run in both server and client mode simultaneously.")
		os.Exit(1)
	}

	if *serverMode {
		log.Info("Starting WebSocket server...")
		websocket.ServeWs()
	} else if *clientMode {
		if *username == "" {
			log.Fatal("Username is required in client mode. Use -user <username>")
		}
		log.Infof("Starting WebSocket client for user %s...\n", *username)
		user := model.User{
			ID:       uuid.NewString(),
			Username: *username,
		}
		conn := websocket.NewConnection(user)
		handler := websocket.NewClientHandler(conn)
		handler.Run()
	} else {
		fmt.Println("Usage: go run main.go [-server | -client -user <username>]")
		os.Exit(1)
	}
	select {}
}
