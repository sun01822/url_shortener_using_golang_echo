package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"url_shortener/internal/redis"

	_ "github.com/joho/godotenv/autoload"

	"url_shortener/internal/database"
)

type Server struct {
	port  int
	db    database.Service
	redis *redis.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:  port,
		db:    database.New(),
		redis: redis.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
