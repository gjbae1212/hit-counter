package main

import (
	"os"
	"testing"

	"github.com/alicebob/miniredis"
)

var (
	mockRedis *miniredis.Miniredis
)

func TestMain(m *testing.M) {
	var err error
	mockRedis, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	code := m.Run()
	mockRedis.Close()
	os.Exit(code)
}
