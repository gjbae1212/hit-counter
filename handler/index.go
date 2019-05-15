package handler

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Index(c echo.Context) error {
	return c.String(http.StatusOK, "index!")
}
