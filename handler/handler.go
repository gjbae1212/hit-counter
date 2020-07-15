package handler

import (
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	task "github.com/gjbae1212/go-async-task"
	websocket "github.com/gjbae1212/go-ws-broadcast"
	"github.com/gjbae1212/hit-counter/counter"
	"github.com/gjbae1212/hit-counter/internal"
	cache "github.com/patrickmn/go-cache"
)

type Handler struct {
	Counter          counter.Counter
	LocalCache       *cache.Cache
	AsyncTask        task.Keeper
	WebSocketBreaker websocket.Breaker
	IndexTemplate    *template.Template
}

// NewHandler creates  handler object.
func NewHandler(redisAddrs []string) (*Handler, error) {
	if len(redisAddrs) == 0 {
		return nil, fmt.Errorf("[err] NewHandler %w", internal.ErrorEmptyParams)
	}

	// create local cache
	localCache := cache.New(24*time.Hour, 10*time.Minute)
	ctr, err := counter.NewCounter(counter.WithRedisOption(redisAddrs))
	if err != nil {
		return nil, fmt.Errorf("[err] NewHandler %w", err)
	}

	// create async task
	asyncTask, err := task.NewAsyncTask(
		task.WithQueueSizeOption(1000),
		task.WithWorkerSizeOption(5),
		task.WithTimeoutOption(20*time.Second),
		task.WithErrorHandlerOption(func(err error) {
			internal.SentryError(err)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("[err] NewHandler %w", err)
	}

	// create websocket breaker
	breaker, err := websocket.NewBreaker(websocket.WithMaxReadLimit(1024),
		websocket.WithMaxMessagePoolLength(500),
		websocket.WithErrorHandlerOption(func(err error) {
			internal.SentryError(err)
		}))
	if err != nil {
		return nil, fmt.Errorf("[err] NewHandler %w", err)
	}

	// template
	indexTemplate, err := template.ParseFiles(filepath.Join(internal.GetRoot(), "view", "index.html"))
	if err != nil {
		return nil, fmt.Errorf("[err] NewHandler %w", err)
	}

	return &Handler{
		LocalCache:       localCache,
		Counter:          ctr,
		AsyncTask:        asyncTask,
		WebSocketBreaker: breaker,
		IndexTemplate:    indexTemplate,
	}, nil
}
