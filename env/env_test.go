package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProject(t *testing.T) {
	assert := assert.New(t)
	p := GetProject()
	assert.Equal(project, p)
}

func TestGetDebug(t *testing.T) {
	assert := assert.New(t)
	d := GetDebug()
	assert.Equal(debugB, d)
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
