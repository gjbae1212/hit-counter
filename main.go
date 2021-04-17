package main // import "github.com/gjbae1212/hit-counter"

import (
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/gjbae1212/hit-counter/internal"

	"path/filepath"

	"github.com/gjbae1212/hit-counter/env"
	"github.com/labstack/echo/v4"
)

var (
	address = flag.String("addr", ":8080", "address")
	tls     = flag.Bool("tls", false, "tls")
)

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	// initialize sentry
	name, _ := os.Hostname()
	if err := internal.InitSentry(env.GetSentryDSN(), env.GetPhase(), env.GetPhase(),
		name, true, env.GetDebug()); err != nil {
		log.Println(err)
	}

	e := echo.New()

	// make options for echo server.
	var opts []Option

	// debug option
	opts = append(opts, WithDebugOption(env.GetDebug()))

	var dir string
	var file string
	if env.GetLogPath() != "" {
		dir, file = filepath.Split(env.GetLogPath())
	}

	// logger option
	logger, err := internal.NewLogger(dir, file)
	if err != nil {
		log.Panic(err)
	}
	opts = append(opts, WithLogger(logger))

	// add middleware
	if err := AddMiddleware(e, opts...); err != nil {
		log.Panic(err)
	}

	// add route
	if err := AddRoute(e, env.GetRedisAddrs()[0]); err != nil {
		log.Panic(err)
	}

	if *tls {
		// start TLS server with let's encrypt certification.
		e.StartAutoTLS(*address)
	} else {
		e.Start(*address)
	}
}
