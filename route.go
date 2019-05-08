package main

import (
	"fmt"

	"encoding/base64"
	"net/http"
	"time"

	allan_util "github.com/gjbae1212/go-module/util"
	"github.com/gjbae1212/hit-counter/handler"
	"github.com/gjbae1212/hit-counter/handler/api"
	"github.com/pkg/errors"
	"github.com/labstack/echo/v4"
)

func AddRoute(e *echo.Echo, redisAddrs []string, cacheSize int) error {
	if e == nil {
		return fmt.Errorf("[Err] AddRoute empty params")
	}

	h, err := handler.NewHandler(redisAddrs, cacheSize)
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
	// HealthCheck
	e.GET("/", h.HealthCheck)

	// GROUP /count
	g1, err := groupApiCount()
	if err != nil {
		return errors.Wrap(err, "[err] AddRoute")
	}
	count := e.Group("/api/count", g1...)
	// badge
	count.GET("/keep/badge.svg", api.KeepCount)

	// GROUP /rank
	return nil
}

func groupApiCount() ([]echo.MiddlewareFunc, error) {
	var chain []echo.MiddlewareFunc
	// Add cookie duration 24 hour.
	cookieFunc := func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hitctx := c.(*handler.HitCounterContext)
			var err error
			cookie := &http.Cookie{}
			if cookie, err = c.Cookie("aid"); err != nil {
				v := fmt.Sprintf("%s-%d", c.RealIP(), time.Now().UnixNano())
				b64 := base64.StdEncoding.EncodeToString([]byte(v))
				cookie = &http.Cookie{
					Name:    "aid",
					Value:   b64,
					Expires: time.Now().Add(24 * time.Hour),
				}
				hitctx.SetCookie(cookie)
			}
			hitctx.Set(cookie.Name, cookie.Value)
			return h(hitctx)
		}
	}
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
			hitctx.Set("id", fmt.Sprintf("%s/%s", host, path))
			return h(hitctx)
		}
	}
	chain = append(chain, cookieFunc, paramFunc)
	return chain, nil
}