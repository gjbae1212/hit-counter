package main

import (
	"testing"
	"time"

	"fmt"
	"net/http/httptest"

	"github.com/gjbae1212/hit-counter/env"
	"github.com/stretchr/testify/assert"
)

func TestLoadSentry(t *testing.T) {
	assert := assert.New(t)

	err := LoadSentry("")
	assert.NoError(err)
}

func TestSendSentry(t *testing.T) {
	assert := assert.New(t)
	SendSentry(nil, nil)

	if env.GetSentryDSN() != "" {
		LoadSentry(env.GetSentryDSN())
		SendSentry(fmt.Errorf("test code error"), httptest.NewRequest("GET", "http://localhost", nil))
	}
	_ = assert
	time.Sleep(2 * time.Second)
}
