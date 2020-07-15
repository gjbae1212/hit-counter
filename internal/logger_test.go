package internal

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	assert := assert.New(t)

	// if log-path don't exist, creating it.
	_, filename, _, ok := runtime.Caller(0)
	assert.True(ok)
	logPath := filepath.Join(path.Dir(filename), "../", "logs")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.Mkdir(logPath, os.ModePerm)
	}

	tests := map[string]struct {
		inputs map[string]string
		isErr  bool
	}{
		"stdout": {inputs: map[string]string{"dir": "", "filename": ""}},
		"file":   {inputs: map[string]string{"dir": filepath.Join(path.Dir("test.log"), "../", "logs"), "filename": "test.log"}},
		"error":  {inputs: map[string]string{"dir": "../empty-folder", "filename": "test.log"}, isErr: true},
	}

	for _, t := range tests {
		logger, err := NewLogger(t.inputs["dir"], t.inputs["filename"])
		assert.Equal(t.isErr, err != nil)
		if err == nil {
			logger.Info("hello world")
		}
	}
}
