package handler

import (
	"fmt"

	"github.com/coocood/freecache"
	"github.com/gjbae1212/go-module/redis"
	"github.com/pkg/errors"
)

type Handler struct {
	localCache *freecache.Cache
	redis      redis.Manager
}

func NewHandler(redisAddrs []string, cacheSize int) (*Handler, error) {
	if len(redisAddrs) == 0 || cacheSize <= 0 {
		return nil, fmt.Errorf("[err] empty params \n")
	}

	localCache := freecache.NewCache(cacheSize)
	redis, err := redis.NewManager(redisAddrs)
	if err != nil {
		return nil, errors.Wrap(err, "[err] initialize redis \n")
	}

	return &Handler{
		localCache: localCache,
		redis:      redis,
	}, nil
}
