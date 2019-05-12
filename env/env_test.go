package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDebug(t *testing.T) {
	assert := assert.New(t)
	d := GetDebug()
	assert.Equal(debug, d)
}

func TestGetSentryDSN(t *testing.T) {
	assert := assert.New(t)
	sd := GetSentryDSN()
	assert.Equal(sentryDsn, sd)
}

func TestGetLogPath(t *testing.T) {
	assert := assert.New(t)
	lp := GetLogPath()
	assert.Equal(logPath, lp)
}

func TestGetRedisAddrs(t *testing.T) {
	assert := assert.New(t)
	addrs := GetRedisAddrs()
	assert.Equal(redisAddrs, addrs)
}
