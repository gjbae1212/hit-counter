package handler

import (
	"github.com/labstack/echo/v4"
)

// IconAll returns icon list.
func (h *Handler) IconAll(c echo.Context) error {
	c.Response().Header().Set("Cache-Control", "max-age=7200, public")
	return c.JSON(200, h.IconsList)
}

// Icon returns icon.svg
func (h *Handler) Icon(c echo.Context) error {
	icon := c.Param("icon")
	svg, ok := h.Icons[icon]
	if !ok {
		return c.NoContent(404)
	} else {
		c.Response().Header().Set("Content-Type", "image/svg+xml")
		c.Response().Header().Set("Cache-Control", "max-age=7200, public")
		return c.String(200, string(svg.Origin))
	}
}
