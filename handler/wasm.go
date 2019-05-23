package handler

import (
	"path/filepath"

	"github.com/gjbae1212/hit-counter/util"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Wasm(c echo.Context) error {
	hctx := c.(*HitCounterContext)
	hctx.Response().Header().Set("Content-Encoding", "gzip")
	return c.File(filepath.Join(util.GetRoot(), "view", "hits.wasm"))
}
