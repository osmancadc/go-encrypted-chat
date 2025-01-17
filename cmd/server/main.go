package main

import (
	"github.com/osmancadc/go-encrypted-chat/pkg/logger"
)

func main() {
	logger := logger.NewLogger("INFO")
	logger.Info(" ---- SECURE CHAT INITIALIZED ----")
}
