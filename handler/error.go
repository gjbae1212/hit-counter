package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// error code 에 대한 응답
func (h *Handler) Error(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.NoContent(code)
}
