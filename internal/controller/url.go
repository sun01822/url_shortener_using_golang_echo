package controller

import (
	"net/http"
	"url_shortener/internal/models"
	"url_shortener/internal/service"

	"github.com/labstack/echo/v4"
)

type UrlController struct {
	svc service.Service
}

func NewUrlController(svc service.Service) *UrlController {
	return &UrlController{
		svc: svc,
	}
}

func (c *UrlController) CreateShortUrl(ctx echo.Context) error {
	var request models.CreateShortUrlRequest
	err := ctx.Bind(&request)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	validationErr := request.Validate()
	if validationErr != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": validationErr.Error()})
	}

	resp, err := c.svc.CreateShortUrl(request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, resp)
}

func (c *UrlController) GetOriginalUrl(ctx echo.Context) error {
	code := ctx.Param("code")
	if code == "" || len(code) < 8 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid short URL code"})
	}

	resp, err := c.svc.GetOriginalUrl(code)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusTemporaryRedirect, map[string]string{"original_url": resp.OriginalUrl})
}

func (c *UrlController) DeleteShortUrl() error {

	return nil
}

func (c *UrlController) UpdateShortUrl() error {

	return nil
}

func (c *UrlController) ListShortUrls() error {

	return nil
}
