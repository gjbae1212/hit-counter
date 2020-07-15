package api_handler

import (
	"fmt"

	"github.com/gjbae1212/hit-counter/handler"
	"github.com/gjbae1212/hit-counter/internal"
)

type Handler struct {
	*handler.Handler
}

// NewHandler creates api handler object.
func NewHandler(h *handler.Handler) (*Handler, error) {
	if h == nil {
		return nil, fmt.Errorf("[err] api handler %w", internal.ErrorEmptyParams)
	}
	return &Handler{Handler: h}, nil
}
