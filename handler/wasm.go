package handler

import (
	"path/filepath"

	"github.com/gjbae1212/hit-counter/internal"
	"github.com/labstack/echo/v4"
)

// Wasm is API for serving wasm file.
func (h *Handler) Wasm(c echo.Context) error {
	hctx := c.(*HitCounterContext)
	hctx.Response().Header().Set("Content-Encoding", "gzip")
	return c.File(filepath.Join(internal.GetRoot(), "view", "hits.wasm"))
}
