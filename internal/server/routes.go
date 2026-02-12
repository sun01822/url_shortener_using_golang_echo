package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	e.GET("/", s.HelloWorldHandler)

	e.GET("/health", s.healthHandler)

	e.GET("/redis-health", s.redisHealthHandler)

	e.GET("/migrate", s.migrateHandler)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) redisHealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.redis.Health())
}

func (s *Server) migrateHandler(c echo.Context) error {
	err := s.db.Migrate()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status": "migration failed",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "migration successful",
	})
}
