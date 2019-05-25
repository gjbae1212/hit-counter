package handler

import (
	"fmt"
	"html/template"

	websocket "github.com/gjbae1212/go-module/websocket"

	"time"

	"path/filepath"

	"github.com/gjbae1212/go-module/async_task"
	"github.com/gjbae1212/hit-counter/counter"
	"github.com/gjbae1212/hit-counter/sentry"
	"github.com/gjbae1212/hit-counter/util"
	cache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

type Handler struct {
	Counter          counter.Counter
	LocalCache       *cache.Cache
	AsyncTask        async_task.Keeper
	WebSocketBreaker websocket.Breaker

	IndexTemplate *template.Template
}

func NewHandler(redisAddrs []string) (*Handler, error) {
	if len(redisAddrs) == 0 {
		return nil, fmt.Errorf("[err] handler empty params \n")
	}

	// create local cache
	localCache := cache.New(24*time.Hour, 10*time.Minute)
	ctr, err := counter.NewCounter(counter.WithRedisOption(redisAddrs))
	if err != nil {
		return nil, errors.Wrap(err, "[err] initialize counter \n")
	}

	// create async task
	asyncTask, err := async_task.NewAsyncTask(
		async_task.WithQueueSizeOption(1000),
		async_task.WithWorkerSizeOption(5),
		async_task.WithTimeoutOption(20*time.Second),
		async_task.WithErrorHandlerOption(func(err error) {
			sentry.SendSentry(err, nil)
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "[err] async task initialize \n")
	}

	// create websocket breaker
	breaker, err := websocket.NewBreaker(websocket.WithMaxReadLimit(1024),
		websocket.WithMaxMessagePoolLength(500),
		websocket.WithErrorHandlerOption(func(err error) {
			sentry.SendSentry(err, nil)
		}))
	if err != nil {
		return nil, errors.Wrap(err, "[err] websocket breaker initialize")
	}

	// template
	indexTemplate, err := template.ParseFiles(filepath.Join(util.GetRoot(), "view", "index.html"))
	if err != nil {
		return nil, errors.Wrap(err, "[err] template initialize")
	}

	return &Handler{
		LocalCache:       localCache,
		Counter:          ctr,
		AsyncTask:        asyncTask,
		WebSocketBreaker: breaker,
		IndexTemplate:    indexTemplate,
	}, nil
}
