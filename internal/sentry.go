package internal

import (
	"errors"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
)

var (
	ErrSentryEmptyParam = errors.New("[err] sentry empty params")
)

// InitSentry is to initialize Sentry setting.
func InitSentry(sentryDSN, environment, release, hostname string, stack, debug bool) error {
	if sentryDSN == "" || environment == "" || release == "" || hostname == "" {
		return ErrSentryEmptyParam
	}

	// if debug is true, it could show detail stack-log.
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDSN,
		Environment:      environment,
		Release:          release,
		ServerName:       hostname,
		AttachStacktrace: stack,
		Debug:            debug,
	}); err != nil {
		return err
	}

	return nil
}

// Error sends an error to Sentry.
func SentryError(err error) {
	if err == nil {
		return
	}
	sentry.CaptureException(err)
}

// ErrorWithEcho sends an error with the Echo of context information to the Sentry.
func SentryErrorWithContext(err error, ctx echo.Context, info map[string]string) {
	if err == nil || ctx == nil {
		return
	}
	var id string
	if info != nil {
		id = info["id"]
	}
	if hub := sentryecho.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{ID: id, IPAddress: ctx.RealIP()})
			scope.SetFingerprint([]string{err.Error()})
			hub.CaptureException(err)
		})
	}
}
