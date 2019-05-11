package sentry

import (
	"net/http"

	allan_sentry "github.com/gjbae1212/go-module/sentry"
	"github.com/getsentry/raven-go"
)

var (
	sentryDSN string
)

func LoadSentry(dsn string) error {
	if dsn == "" {
		return nil
	}
	sentryDSN = dsn
	return allan_sentry.Load(dsn, 0)
}

func SendSentry(err error, r *http.Request) {
	if sentryDSN == "" || err == nil {
		return
	}
	var packet *raven.Packet
	if r != nil {
		packet = allan_sentry.MakePacketWithRequest(err, r, true)
	} else {
		packet = allan_sentry.MakePacket(err, true)
	}
	allan_sentry.Raise(packet, nil)
}
