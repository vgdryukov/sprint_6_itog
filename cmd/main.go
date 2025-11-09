package main

import (
	"log"
	"os"

	"internal/server"
)

func main() {
	// Создаем логгер
	logger := log.New(os.Stdout, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Создаем сервер
	srv := server.NewServer(logger)

	// Запускаем сервер
	if err := srv.Start(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
