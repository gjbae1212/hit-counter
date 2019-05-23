package handler

import (
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/gjbae1212/hit-counter/util"
)

func (h *Handler) Index(c echo.Context) error {
	buf, err := ioutil.ReadFile(filepath.Join(util.GetRoot(), "view", "index.html"))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, buf)
}
