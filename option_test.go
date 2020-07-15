package main

import (
	"testing"

	glog "github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestWithDebugOption(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()

	tests := map[string]struct {
		input  bool
		output bool
	}{
		"false": {},
		"true":  {},
	}

	for _, t := range tests {
		opt := WithDebugOption(t.input)
		opt.apply(e)
		assert.Equal(e.Debug, t.output)
	}
}

func TestWithLogger(t *testing.T) {
	assert := assert.New(t)
	e := echo.New()

	tests := map[string]struct {
		input  *glog.Logger
		output *glog.Logger
	}{
		"nil":   {input: nil, output: nil},
		"exist": {},
	}

	for _, t := range tests {
		opt := WithLogger(t.input)
		opt.apply(e)
		assert.Equal(e.Logger, t.input)
	}
}
