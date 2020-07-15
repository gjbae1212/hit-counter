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
	forceHTTPS bool
	redisAddrs []string
	cacheSize  int
	phase      string
)

func init() {
	var err error
	if os.Getenv("DEBUG") != "" {
		debug, err = strconv.ParseBool(os.Getenv("DEBUG"))
		if err != nil {
			log.Panic(err)
		}
	}
	if os.Getenv("FORCE_HTTPS") != "" {
		forceHTTPS, err = strconv.ParseBool(os.Getenv("FORCE_HTTPS"))
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

// GetDebug returns DEBUG global environment.
func GetDebug() bool {
	return debug
}

// GetDebug returns LOG_PATH global environment.
func GetLogPath() string {
	return logPath
}

// GetSentryDSN returns SENTRY_DSN global environment.
func GetSentryDSN() string {
	return sentryDsn
}

// GetRedisAddrs returns REDIS_ADDRS global environment.
func GetRedisAddrs() []string {
	return redisAddrs
}

// GetForceHTTPS returns FORCE_HTTPS global environment.
func GetForceHTTPS() bool {
	return forceHTTPS
}

// GetPhase returns PHASE global environment.
func GetPhase() string {
	return phase
}
