package handler

import (
	"fmt"

	websocket "github.com/gjbae1212/go-ws-broadcast"
	"github.com/labstack/echo/v4"
)

// WebSocket is API for websocket.
func (h *Handler) WebSocket(c echo.Context) error {
	ws, err := websocket.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return fmt.Errorf("[err] WebSocket API %w", err)
	}

	// register websocket to breaker.
	if _, err := h.WebSocketBreaker.Register(ws); err != nil {
		return fmt.Errorf("[err] WebSocket API %w", err)
	}
	return nil
}
