package server

import (
	"log"
	"net/http"
	"time"

	"sprint_6_itog/internal/handlers"
	"sprint_6_itog/internal/service"
)

type Server struct {
	logger *log.Logger
	server *http.Server
}

func NewServer(logger *log.Logger) *Server {
	// Создаем сервис и хендлеры
	svc := service.NewService()
	handler := handlers.NewHandler(svc)

	// Создаем роутер
	router := http.NewServeMux()
	router.HandleFunc("/", handler.RootHandler)
	router.HandleFunc("/upload", handler.UploadHandler)

	// Создаем HTTP сервер
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		logger: logger,
		server: httpServer,
	}
}

// Start запускает сервер
func (s *Server) Start() error {
	s.logger.Printf("Starting server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}
