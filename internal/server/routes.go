package server

import (
	"net/http"
	middleware2 "url_shortener/internal/middleware"

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

	e.GET("/db-health", s.healthHandler)

	//e.GET("/redis-health", s.redisHealthHandler)

	e.GET("/migrate", s.migrateHandler, middleware2.CustomAuth())
	e.GET("/clear", s.clearHandler, middleware2.CustomAuth())

	u := e.Group("/api/v1")
	{
		u.POST("/shorten", s.controller.CreateShortUrl)
		u.GET("/shorten/:code", s.controller.GetOriginalUrl)
	}

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.dbSvc.Health())
}

//func (s *Server) redisHealthHandler(c echo.Context) error {
//	return c.JSON(http.StatusOK, s.redis.Health())
//}

func (s *Server) migrateHandler(c echo.Context) error {
	err := s.dbSvc.Migrate()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Migration successful"})
}

func (s *Server) clearHandler(c echo.Context) error {
	err := s.dbSvc.Clear()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Clear all url successful"})
}
