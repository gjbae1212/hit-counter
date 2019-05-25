package api_handler

import (
	"net/http"

	"github.com/gjbae1212/go-badge"

	"time"

	"fmt"

	"context"

	"github.com/cespare/xxhash"
	allan_util "github.com/gjbae1212/go-module/util"
	"github.com/gjbae1212/hit-counter/counter"
	"github.com/gjbae1212/hit-counter/handler"
	"github.com/gjbae1212/hit-counter/sentry"
	"github.com/labstack/echo/v4"
)

var (
	badgeFormat   = " %d / %d "
	countIdFormat = "%s%s"

	domainGroup = "domain"
)

type RankTask struct {
	Counter   counter.Counter
	Domain    string
	Path      string
	CreatedAt time.Time
}

func (task *RankTask) Process(ctx context.Context) error {
	// If a domain is 'github.com', it is calculating ranks.
	if task.Domain == "github.com" && task.Path != "" {
		if _, err := task.Counter.IncreaseRankOfDaily(task.Domain, task.Path, task.CreatedAt); err != nil {
			return err
		}
		if _, err := task.Counter.IncreaseRankOfTotal(task.Domain, task.Path); err != nil {
			return err
		}
	}

	// Calculate ranks for a domain of daily and total
	if _, err := task.Counter.IncreaseRankOfDaily(domainGroup, task.Domain, task.CreatedAt); err != nil {
		return err
	}
	if _, err := task.Counter.IncreaseRankOfTotal(domainGroup, task.Domain); err != nil {
		return err
	}

	return nil
}

type WebSocketMessage struct {
	Payload []byte
}

func (wsm *WebSocketMessage) GetMessage() []byte {
	return wsm.Payload
}

func (h *Handler) IncrCount(c echo.Context) error {
	hctx := c.(*handler.HitCounterContext)
	if hctx.Get("ckid") == nil || hctx.Get("host") == nil || hctx.Get("path") == nil {
		return fmt.Errorf("[err] IncrCount empty params")
	}
	cookie := hctx.Get("ckid").(string)
	host := hctx.Get("host").(string)
	path := hctx.Get("path").(string)
	_ = cookie
	id := fmt.Sprintf(countIdFormat, host, path)
	ip := c.RealIP()
	userAgent := c.Request().UserAgent()

	// If a ingress specified ip is exceeded more than 30 per 5 seconds, it might possibly abusing.
	// so it must be limited.
	v, ok := h.LocalCache.Get(ip)
	if v != nil && v.(int64) > 20 {
		daily, total, err := h.Counter.GetHitAll(id, time.Now())
		if err != nil {
			return err
		}
		return h.returnCount(hctx, daily, total)
	}
	if !ok {
		h.LocalCache.Set(ip, int64(1), 5*time.Second)
	} else {
		if _, err := h.LocalCache.IncrementInt64(ip, 1); err != nil {
			return err
		}
	}

	// It would limit for count a hit when a user is accessing more than 1 per second.
	temporaryId := fmt.Sprintf("%d", xxhash.Sum64String(fmt.Sprintf("%s-%s", ip, userAgent)))
	if _, ok := h.LocalCache.Get(temporaryId); ok {
		daily, total, err := h.Counter.GetHitAll(id, time.Now())
		if err != nil {
			return err
		}
		return h.returnCount(hctx, daily, total)
	}
	h.LocalCache.Set(temporaryId, int64(1), 1*time.Second)

	daily, err := h.Counter.IncreaseHitOfDaily(id, time.Now())
	if err != nil {
		return err
	}

	total, err := h.Counter.IncreaseHitOfTotal(id)
	if err != nil {
		return err
	}

	// Calculating ranks of daily and total is asynchronously working.
	timectx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.AsyncTask.AddTask(timectx, &RankTask{
		Counter:   h.Counter,
		Domain:    host,
		Path:      path,
		CreatedAt: time.Now(),
	}); err != nil { // Possibly send an error to the sentry. And it is not returned a error.
		sentry.SendSentry(err, nil)
	}

	// Broadcast message to users to which connected
	h.WebSocketBreaker.BroadCast(&WebSocketMessage{
		Payload: []byte(fmt.Sprintf("[%s] %s", allan_util.TimeToString(time.Now()), id))},
	)
	return h.returnCount(hctx, daily, total)
}

func (h *Handler) KeepCount(c echo.Context) error {
	hctx := c.(*handler.HitCounterContext)
	if hctx.Get("ckid") == nil || hctx.Get("host") == nil || hctx.Get("path") == nil {
		return fmt.Errorf("[err] KeepCount empty params")
	}
	host := hctx.Get("host").(string)
	path := hctx.Get("path").(string)
	cookie := hctx.Get("ckid").(string)
	_ = cookie
	id := fmt.Sprintf(countIdFormat, host, path)
	daily, total, err := h.Counter.GetHitAll(id, time.Now())
	if err != nil {
		return err
	}

	return h.returnCount(hctx, daily, total)
}

func (h *Handler) returnCount(hctx *handler.HitCounterContext, daily, total *counter.Score) error {
	dailyCount := int64(0)
	totalCount := int64(0)
	if daily != nil {
		dailyCount = daily.Value
	}
	if total != nil {
		totalCount = total.Value
	}

	text := fmt.Sprintf(badgeFormat, dailyCount, totalCount)
	badge, err := badge.RenderBytes("hits", text, "#79c83d")
	if err != nil {
		return err
	}

	hctx.Response().Header().Set("Content-Type", "image/svg+xml")
	hctx.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	hctx.Response().Header().Set("Pragma", "no-cache")
	hctx.Response().Header().Set("Expires", "0")
	return hctx.String(http.StatusOK, string(badge))
}
