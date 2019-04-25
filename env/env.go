package env

import (
	"log"
	"os"
	"strconv"
)

var (
	project   = os.Getenv("PROJECT")
	debug     = os.Getenv("DEBUG")
	logPath   = os.Getenv("LOG_PATH")
	sentryDsn = os.Getenv("SENTRY_DSN")
)

var (
	debugB bool
)

func init() {
	debugB = true
	if debug != "" {
		var err error
		debugB, err = strconv.ParseBool(debug)
		if err != nil {
			log.Panic(err)
		}
	}
}

func GetProject() string {
	return project
}

func GetDebug() bool {
	return debugB
}

func GetLogPath() string {
	return logPath
}

func GetSentryDSN() string {
	return sentryDsn
}
