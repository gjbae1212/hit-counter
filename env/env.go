package env

import (
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	logPath   = os.Getenv("LOG_PATH")
	sentryDsn = os.Getenv("SENTRY_DSN")
)

var (
	debug      bool
	redisAddrs []string
	cacheSize  int
)

func init() {
	var err error
	debug = true
	if os.Getenv("DEBUG") != "" {
		debug, err = strconv.ParseBool(os.Getenv("DEBUG"))
		if err != nil {
			log.Panic(err)
		}
	}

	if os.Getenv("REDIS_ADDRS") != "" {
		seps := strings.Split(os.Getenv("REDIS_ADDRS"), ",")
		for _, sep := range seps {
			redisAddrs = append(redisAddrs, strings.TrimSpace(sep))
		}
	}

}

func GetDebug() bool {
	return debug
}

func GetLogPath() string {
	return logPath
}

func GetSentryDSN() string {
	return sentryDsn
}

func GetRedisAddrs() []string {
	return redisAddrs
}
