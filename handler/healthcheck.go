package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "health check!")
}
