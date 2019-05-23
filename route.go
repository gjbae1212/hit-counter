package main

import (
	"fmt"

	"net/http"
	allan_util "github.com/gjbae1212/go-module/util"
	"github.com/gjbae1212/hit-counter/handler"
	"github.com/gjbae1212/hit-counter/handler/api"
	"github.com/pkg/errors"
	"github.com/labstack/echo/v4"
)

func AddRoute(e *echo.Echo, redisAddrs []string) error {
	if e == nil {
		return fmt.Errorf("[Err] AddRoute empty params")
	}

	h, err := handler.NewHandler(redisAddrs)
	if err != nil {
		return errors.Wrap(err, "[err] AddRoute")
	}

	api, err := api_handler.NewHandler(h)
	if err != nil {
		return errors.Wrap(err, "[err] AddRoute")
	}

	// error handler
	e.HTTPErrorHandler = h.Error
	// static
	e.Static("/", "public")

	// wasm
	e.GET("/hits.wasm", h.Wasm)

	// main
	e.GET("/", h.Index)

	// health check
	e.GET("/healthcheck", h.HealthCheck)

	// group /count
	g1, err := groupApiCount()
	if err != nil {
		return errors.Wrap(err, "[err] AddRoute")
	}
	count := e.Group("/api/count", g1...)
	// badge
	count.GET("/keep/badge.svg", api.KeepCount)
	count.GET("/incr/badge.svg", api.IncrCount)

	// GROUP /rank
	return nil
}

func groupApiCount() ([]echo.MiddlewareFunc, error) {
	var chain []echo.MiddlewareFunc
	// Add param
	paramFunc := func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hitctx := c.(*handler.HitCounterContext)
			// url validation check
			url := hitctx.QueryParam("url")
			if url == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Not Found URL Query String")
			}

			schema, host, _, path, _, _, err := allan_util.ParseURL(url)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid URL Query String %s", url))
			}

			if !allan_util.StringInSlice(schema, []string{"http", "https"}) {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Not Support Schema %s", schema))
			}
			hitctx.Set("host", host)
			hitctx.Set("path", path)
			return h(hitctx)
		}
	}
	chain = append(chain, paramFunc)
	return chain, nil
}
