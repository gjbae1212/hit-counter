package handler

import (
	"github.com/pkg/errors"

	websocket "github.com/gjbae1212/go-module/websocket"
	"github.com/labstack/echo/v4"
)

func (h *Handler) WebSocket(c echo.Context) error {
	ws, err := websocket.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return errors.Wrap(err, "[err] websocket")
	}

	if err := h.WebSocketBreaker.Register(ws); err != nil {
		return errors.Wrap(err, "[err] websocket register")
	}
	return nil
}
