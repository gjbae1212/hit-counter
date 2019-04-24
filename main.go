package main // import "github.com/gjbae1212/hit-counter"

import (
	"flag"
	"log"

	"runtime"
	echo "github.com/labstack/echo/v4"
)

var (
	address = flag.String("addr", ":8080", "address")
	tls     = flag.Bool("tls", false, "tls")
)

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	e := echo.New()
	
	if *tls {
		// If it use to `let's encrypt`
		e.StartAutoTLS(*address)
	} else {
		e.Start(*address)
	}
}
