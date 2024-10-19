package server

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"new_project/internal/database"
)

type Server struct {
	port   int
	db     database.Service
	logger *slog.Logger
	wg     sync.WaitGroup
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	NewServer := &Server{
		port:   port,
		db:     database.New(),
		logger: logger,
	}

	// Declare Server config
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(),
		//IdleTimeout:  time.Minute,
		//ReadTimeout:  10 * time.Second,
		//WriteTimeout: 30 * time.Second,
		//debugging
		IdleTimeout:  5 * time.Minute,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

	return server
}
