package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gjbae1212/hit-counter/handler"
	api_handler "github.com/gjbae1212/hit-counter/handler/api"
	"github.com/gjbae1212/hit-counter/internal"
	"github.com/labstack/echo/v4"
)

func AddRoute(e *echo.Echo, redisAddrs []string) error {
	if e == nil {
		return fmt.Errorf("[Err] AddRoute empty params")
	}

	h, err := handler.NewHandler(redisAddrs)
	if err != nil {
		return fmt.Errorf("[err] AddRoute %w", err)
	}

	api, err := api_handler.NewHandler(h)
	if err != nil {
		return fmt.Errorf("[err] AddRoute %w", err)
	}

	// error handler
	e.HTTPErrorHandler = h.Error
	// static
	e.Static("/", "public")

	// wasm
	e.GET("/hits.wasm", h.Wasm)

	// websocket
	e.GET("/ws", h.WebSocket)

	// main
	e.GET("/", h.Index)

	// icon
	e.GET("/icon/all.json", h.IconAll)
	e.GET("/icon/:icon", h.Icon)

	// health check
	e.GET("/healthcheck", h.HealthCheck)

	// group /api/count
	g1, err := groupApiCount()
	if err != nil {
		return fmt.Errorf("[err] AddRoute %w", err)
	}
	count := e.Group("/api/count", g1...)
	// badge
	count.GET("/keep/badge.svg", api.KeepCount)
	count.GET("/incr/badge.svg", api.IncrCount)

	// graph
	count.GET("/graph/dailyhits.svg", api.DailyHitsInRecently)

	// group /api/rank
	g2, err := groupApiRank()
	if err != nil {
		return fmt.Errorf("[err] AddRoute %w", err)
	}
	rank := e.Group("/api/rank", g2...)
	_ = rank

	return nil
}

func groupApiCount() ([]echo.MiddlewareFunc, error) {
	var chain []echo.MiddlewareFunc
	// Add param
	paramFunc := func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hitctx := c.(*handler.HitCounterContext)

			// check a url is invalid or not.
			url := hitctx.QueryParam("url")
			if url == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Not Found URL Query String")
			}

			schema, host, _, path, _, _, err := internal.ParseURL(url)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid URL Query String %s", url))
			}

			if !internal.StringInSlice(schema, []string{"http", "https"}) {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Not Support Schema %s", schema))
			}

			// extract required parameters
			title := hitctx.QueryParam("title")
			titleBg := hitctx.QueryParam("title_bg")
			countBg := hitctx.QueryParam("count_bg")
			edgeFlat, _ := strconv.ParseBool(hitctx.QueryParam("edge_flat"))
			icon := hitctx.QueryParam("icon")
			iconColor := hitctx.QueryParam("icon_color")

			// insert params to context.
			hitctx.Set("host", host)
			hitctx.Set("path", path)
			hitctx.Set("title", title)
			hitctx.Set("title_bg", titleBg)
			hitctx.Set("count_bg", countBg)
			hitctx.Set("edge_flat", edgeFlat)
			hitctx.Set("icon", icon)
			hitctx.Set("icon_color", iconColor)

			return h(hitctx)
		}
	}
	chain = append(chain, paramFunc)
	return chain, nil
}

func groupApiRank() ([]echo.MiddlewareFunc, error) {
	var chain []echo.MiddlewareFunc
	return chain, nil
}
