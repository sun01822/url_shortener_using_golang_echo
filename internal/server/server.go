package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"url_shortener/internal/controller"
	"url_shortener/internal/redis"
	"url_shortener/internal/repository"
	"url_shortener/internal/service"

	"url_shortener/internal/database"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

type Server struct {
	port       int
	db         *gorm.DB
	dbSvc      *database.Service
	redis      *redis.Service
	controller *controller.UrlController
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:  port,
		db:    database.New(),
		redis: redis.New(),
	}

	NewServer.dbSvc = database.GetService(NewServer.db)

	newDB := repository.NewUrlRepository(NewServer.db)
	newService := service.NewUrlService(newDB)
	newController := controller.NewUrlController(newService)
	NewServer.controller = newController

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
