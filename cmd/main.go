package main

import (
	"log"
	"os"

	"sprint_6_itog/internal/server"
)

func main() {
	// Открываем файл для логирования
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл логов: %v", err)
	}
	defer logFile.Close()

	// Создаем логгер, который пишет в файл
	logger := log.New(logFile, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile)

	srv := server.NewServer(logger)

	if err := srv.Start(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
