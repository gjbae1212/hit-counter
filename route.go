package main

import (
	"fmt"

	"github.com/gjbae1212/hit-counter/handler"
	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func AddRoute(e *echo.Echo, redisAddrs []string, cacheSize int) error {
	if e == nil {
		return fmt.Errorf("[Err] AddRoute empty params")
	}

	h, err := handler.NewHandler(redisAddrs, cacheSize)
	if err != nil {
		return errors.Wrap(err, "[err] AddRoute")
	}

	// error handler
	e.HTTPErrorHandler = h.Error

	// static
	e.Static("/", "public")

	// GET
	e.GET("/", h.HealthCheck)

	// POST
	return nil
}
