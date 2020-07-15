package main

import (
	"github.com/labstack/echo/v4"
	glog "github.com/labstack/gommon/log"
)

// Option is an interface for dependency injection in Echo.
type Option interface {
	apply(e *echo.Echo)
}

type OptionFunc func(e *echo.Echo)

func (f OptionFunc) apply(e *echo.Echo) { f(e) }

// WithDebugOption returns a function which sets debug variable in echo server.
func WithDebugOption(debug bool) OptionFunc {
	return func(e *echo.Echo) {
		e.Debug = debug
	}
}

// WithLogger returns a function which sets logger variable in echo server.
func WithLogger(logger *glog.Logger) OptionFunc {
	return func(e *echo.Echo) {
		e.Logger = logger
	}
}
