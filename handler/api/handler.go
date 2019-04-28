package api_handler

import (
	"fmt"

	"github.com/gjbae1212/hit-counter/handler"
)

type Handler struct {
	*handler.Handler
}

func NewHandler(h *handler.Handler) (*Handler, error) {
	if h == nil {
		return nil, fmt.Errorf("[err] api handler empty params\n")
	}
	return &Handler{Handler: h}, nil
}