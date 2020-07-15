package internal

import (
	"os"
	"path/filepath"

	"github.com/labstack/gommon/log"
)

// NewLogger is to create a logger object for Echo platform.
func NewLogger(dir, filename string) (*log.Logger, error) {
	logger := log.New("")
	logger.SetHeader("{\"time\":\"${time_rfc3339}\", \"level\":\"${level}\"}")
	fpath := filepath.Join(dir, filename)
	if fpath == "" {
		logger.SetOutput(os.Stdout)
	} else {
		f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		logger.SetOutput(f)
	}
	return logger, nil
}
