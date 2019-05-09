package handler

import (
	"fmt"

	"github.com/coocood/freecache"
	"github.com/gjbae1212/hit-counter/counter"
	"github.com/pkg/errors"
)

type Handler struct {
	Counter    counter.Counter
	LocalCache *freecache.Cache
}

func NewHandler(redisAddrs []string, cacheSize int) (*Handler, error) {
	if len(redisAddrs) == 0 || cacheSize <= 0 {
		return nil, fmt.Errorf("[err] handler empty params \n")
	}

	localCache := freecache.NewCache(cacheSize)
	ctr, err := counter.NewCounter(counter.WithRedisOption(redisAddrs))
	if err != nil {
		return nil, errors.Wrap(err, "[err] initialize counter \n")
	}

	return &Handler{
		LocalCache: localCache,
		Counter:    ctr,
	}, nil
}
