package handler

import (
	"context"
	"github.com/labstack/echo/v4"
)

type (
	HitCounterContext struct {
		echo.Context
	}
)

func (c *HitCounterContext) GetContext() context.Context {
	return c.Request().Context()
}

func (c *HitCounterContext) SetContext(ctx context.Context) {
	c.SetRequest(c.Request().WithContext(ctx))
}

func (c *HitCounterContext) WithContext(key, val interface{}) {
	ctx := c.GetContext()
	c.SetContext(context.WithValue(ctx, key, val))
}

func (c *HitCounterContext) ValueContext(key interface{}) interface{} {
	return c.GetContext().Value(key)
}

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
