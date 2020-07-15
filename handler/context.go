package handler

import (
	"context"

	"github.com/labstack/echo/v4"
)

type (
	// It's custom context.
	HitCounterContext struct {
		echo.Context
	}
)

// GetContext returns a context in request.
func (c *HitCounterContext) GetContext() context.Context {
	return c.Request().Context()
}

// SetContext sets a context to request.
func (c *HitCounterContext) SetContext(ctx context.Context) {
	c.SetRequest(c.Request().WithContext(ctx))
}

// WithContext set a context with new value to request.
func (c *HitCounterContext) WithContext(key, val interface{}) {
	ctx := c.GetContext()
	c.SetContext(context.WithValue(ctx, key, val))
}

// ValueContext returns values in request context.
func (c *HitCounterContext) ValueContext(key interface{}) interface{} {
	return c.GetContext().Value(key)
}

// ExtraLog returns log struct.
func (c *HitCounterContext) ExtraLog() map[string]interface{} {
	return map[string]interface{}{
		"host":       c.Request().Host,
		"ip":         c.RealIP(),
		"uri":        c.Request().RequestURI,
		"method":     c.Request().Method,
		"referer":    c.Request().Referer(),
		"user-agent": c.Request().UserAgent(),
	}
}
