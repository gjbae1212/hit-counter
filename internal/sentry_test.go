package internal

import (
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gjbae1212/hit-counter/env"
	"github.com/labstack/echo/v4"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestInitSentry(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		inputSentryDSN   string
		inputEnvironment string
		inputRelease     string
		inputHostname    string
		inputStack       bool
		inputDebug       bool
		outputError      bool
	}{
		"empty":   {outputError: true},
		"success": {inputSentryDSN: env.GetSentryDSN(), inputEnvironment: "local", inputRelease: "test", inputHostname: "localhost"},
	}

	for _, t := range tests {
		//err := InitSentry(t.inputSentryDSN, t.inputEnvironment, t.inputRelease, t.inputHostname, false, false)
		//assert.Equal(t.outputError, err != nil)
		_ = t
		_ = assert
	}
}

func TestSentryError(t *testing.T) {
	assert := assert.New(t)

	InitSentry(env.GetSentryDSN(), "local", "test", "localhost", true, true)

	tests := map[string]struct {
		inputErr error
	}{
		"success": {inputErr: fmt.Errorf("[err] test default")},
	}

	for _, t := range tests {
		SentryError(t.inputErr)
		sentry.Flush(5 * time.Second)
	}
	time.Sleep(2 * time.Second)
	_ = assert
}

func TestSentryErrorWithContext(t *testing.T) {
	assert := assert.New(t)

	InitSentry(env.GetSentryDSN(), "local", "test", "localhost", true, true)

	e := echo.New()

	tests := map[string]struct {
		inputErr     error
		inputContext echo.Context
	}{
		"success": {inputErr: fmt.Errorf("[err] test context"), inputContext: e.NewContext(nil, nil)},
	}

	for _, t := range tests {
		SentryErrorWithContext(t.inputErr, t.inputContext, nil)
		sentry.Flush(5 * time.Second)
	}
	time.Sleep(2 * time.Second)
	_ = assert
}
