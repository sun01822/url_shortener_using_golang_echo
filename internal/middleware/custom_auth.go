package middleware

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func CustomAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			appKeyHeader := os.Getenv("APP_KEY_HEADER")
			appKey := os.Getenv("APP_KEY")
			headers := c.Request().Header
			key := headers.Get(appKeyHeader)

			if appKey == key {
				return next(c)
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"message": "Has no permission to operate",
			})
		}
	}
}
