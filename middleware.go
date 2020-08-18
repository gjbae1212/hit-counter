package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/gjbae1212/hit-counter/env"
	"github.com/gjbae1212/hit-counter/handler"
	"github.com/gjbae1212/hit-counter/internal"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// AddMiddleware adds middlewares to echo server.
func AddMiddleware(e *echo.Echo, opts ...Option) error {
	if e == nil {
		return fmt.Errorf("[err] AddMiddleware %w", internal.ErrorEmptyParams)
	}

	o := []Option{WithDebugOption(true)}
	o = append(o, opts...)
	for _, opt := range o {
		opt.apply(e)
	}

	e.HideBanner = true
	e.HidePort = true

	// read timeout will wait until read to request body
	e.Server.ReadTimeout = 10 * time.Second

	// write timeout will wait until from read request body to write response
	e.Server.WriteTimeout = 10 * time.Second

	// pre chain middleware
	prechain, err := middlewarePreChain()
	if err != nil {
		return err
	}
	e.Use(prechain...)

	// main chain middleware
	mainchain, err := middlewareChain()
	if err != nil {
		return err
	}
	e.Use(mainchain...)
	return nil
}

func middlewarePreChain() ([]echo.MiddlewareFunc, error) {
	var chain []echo.MiddlewareFunc

	// set sentry middleware. if this middleware will catch a panic error, delivering it to upper middleware.
	chain = append(chain, sentryecho.New(sentryecho.Options{Repanic: true}))

	// custom context
	if env.GetForceHTTPS() {
		// Apply HSTS
		chain = append(chain, func(h echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Response().Header().Set("Strict-Transport-Security",
					"max-age=2592000; includeSubdomains; preload")
				return h(c)
			}
		})
		// Redirect Https
		chain = append(chain, middleware.HTTPSRedirect())
	}
	chain = append(chain, middleware.RemoveTrailingSlash())
	chain = append(chain, middleware.NonWWWRedirect())
	chain = append(chain, middleware.Rewrite(map[string]string{
		"/static/*": "/public/$1",
	}))

	// Add custom context
	chain = append(chain, func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hitctx := &handler.HitCounterContext{c}

			// set start time
			hitctx.WithContext("start_time", time.Now())

			// set deadline
			timeout := 15 * time.Second

			ctx, cancel := context.WithTimeout(hitctx.GetContext(), timeout)
			defer cancel()
			hitctx.SetContext(ctx)

			// set log
			extraLog := hitctx.ExtraLog()
			hitctx.WithContext("extra_log", extraLog)
			return h(hitctx)
		}
	})

	// Add cookie duration 24 hour.
	chain = append(chain, func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hitctx := c.(*handler.HitCounterContext)
			var err error
			cookie := &http.Cookie{}
			if cookie, err = c.Cookie("ckid"); err != nil {
				v := fmt.Sprintf("%s-%d", c.RealIP(), time.Now().UnixNano())
				b64 := base64.StdEncoding.EncodeToString([]byte(v))
				cookie = &http.Cookie{
					Name:     "ckid",
					Value:    b64,
					Expires:  time.Now().Add(24 * time.Hour),
					Path:     "/",
					HttpOnly: true,
				}
				hitctx.SetCookie(cookie)
			}
			hitctx.Set(cookie.Name, cookie.Value)
			return h(hitctx)
		}
	})
	return chain, nil
}

func middlewareChain() ([]echo.MiddlewareFunc, error) {
	var chain []echo.MiddlewareFunc

	// main middleware
	m := func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// recover
			defer func() {
				if r := recover(); r != nil {
					// send sentry
					internal.SentryErrorWithContext(r.(error), c, nil)

					extraLog := c.(*handler.HitCounterContext).ValueContext("extra_log").(map[string]interface{})
					extraLog["status"] = http.StatusInternalServerError
					extraLog["error"] = fmt.Sprintf("%v\n", r)
					c.Logger().Errorj(extraLog)
					// send 500 error
					c.NoContent(http.StatusInternalServerError)
				}
			}()

			hitctx := c.(*handler.HitCounterContext)
			start := hitctx.ValueContext("start_time").(time.Time)

			// main handler process
			err := h(hitctx)
			stop := time.Now()
			if err != nil {
				code := http.StatusInternalServerError
				if he, ok := err.(*echo.HTTPError); ok {
					code = he.Code
				} else if hitctx.Response().Status >= 400 {
					code = hitctx.Response().Status
				}

				extraLog := hitctx.ValueContext("extra_log").(map[string]interface{})
				extraLog["status"] = code
				extraLog["error"] = fmt.Sprintf("%v\n", err)
				if code >= http.StatusInternalServerError {
					// send sentry
					internal.SentryErrorWithContext(err, c, nil)

					rest := stop.Sub(start)
					extraLog["latency"] = strconv.FormatInt(int64(rest), 10)
					extraLog["latency_human"] = rest.String()
				}
				hitctx.Logger().Errorj(extraLog)
				return err
			}
			extraLog := hitctx.ValueContext("extra_log").(map[string]interface{})
			if extraLog["uri"] != "/healthcheck" {
				extraLog["status"] = hitctx.Response().Status
				rest := stop.Sub(start)
				extraLog["latency"] = strconv.FormatInt(int64(rest), 10)
				extraLog["latency_human"] = rest.String()
				//hitctx.Logger().Infoj(extraLog)
			}
			return nil
		}
	}
	chain = append(chain, m)
	return chain, nil
}
