package main

import (
	"net/http"

	"github.com/gjbae1212/go-module/sentry"
)

var (
	sentryDSN string
)

func LoadSentry(dsn string) error {
	if dsn == "" {
		return nil
	}
	sentryDSN = dsn
	return sentry.Load(dsn, 0)
}

func SendSentry(err error, r *http.Request) {
	if sentryDSN == "" {
		return
	}
	packet := sentry.MakePacketWithRequest(err, r, true)
	sentry.Raise(packet, nil)
}
