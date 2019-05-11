package main // import "github.com/gjbae1212/hit-counter"

import (
	"flag"
	"log"
	"runtime"

	"path/filepath"

	"github.com/gjbae1212/go-module/logger"
	"github.com/gjbae1212/hit-counter/env"
	"github.com/gjbae1212/hit-counter/sentry"
	"github.com/labstack/echo/v4"
)

var (
	address = flag.String("addr", ":8080", "address")
	tls     = flag.Bool("tls", false, "tls")
)

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	// If the sentry is possibly loaded
	sentry.LoadSentry(env.GetSentryDSN())

	e := echo.New()

	// options
	var opts []Option
	opts = append(opts, WithDebugOption(env.GetDebug()))

	dir := ""
	file := ""
	if env.GetLogPath() != "" {
		dir, file = filepath.Split(env.GetLogPath())
	}
	customLogger, err := logger.NewLogger(dir, file)
	if err != nil {
		log.Panic(err)
	}
	opts = append(opts, WithLoggerOption(customLogger))

	// add middleware
	if err := AddMiddleware(e, opts...); err != nil {
		log.Panic(err)
	}

	// add route
	if err := AddRoute(e, env.GetRedisAddrs(), env.GetCacheSize()); err != nil {
		log.Panic(err)
	}

	if *tls {
		// If it use to `let's encrypt`
		e.StartAutoTLS(*address)
	} else {
		e.Start(*address)
	}
}
