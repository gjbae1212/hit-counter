package api_handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cespare/xxhash"
	"github.com/gjbae1212/hit-counter/counter"
	"github.com/gjbae1212/hit-counter/handler"
	"github.com/gjbae1212/hit-counter/internal"
	"github.com/labstack/echo/v4"
	"github.com/wcharczuk/go-chart"
)

var (
	badgeFormat   = " %d / %d "
	countIdFormat = "%s%s"
)

// IncrCount is API, which it's to increase page count.
func (h *Handler) IncrCount(c echo.Context) error {
	hctx := c.(*handler.HitCounterContext)
	if hctx.Get("ckid") == nil || hctx.Get("host") == nil || hctx.Get("path") == nil ||
		hctx.Get("title") == nil || hctx.Get("title_bg") == nil || hctx.Get("count_bg") == nil ||
		hctx.Get("edge_flat") == nil || hctx.Get("icon") == nil || hctx.Get("icon_color") == nil {
		return fmt.Errorf("[err] IncrCount empty params")
	}
	cookie := hctx.Get("ckid").(string)
	host := hctx.Get("host").(string)
	path := hctx.Get("path").(string)
	title := hctx.Get("title").(string)
	titleBg := hctx.Get("title_bg").(string)
	countBg := hctx.Get("count_bg").(string)
	edgeType := hctx.Get("edge_flat").(bool)
	icon := hctx.Get("icon").(string)
	iconColor := hctx.Get("icon_color").(string)

	_ = cookie
	id := fmt.Sprintf(countIdFormat, host, path)
	ip := c.RealIP()
	userAgent := c.Request().UserAgent()

	// If a ingress specified ip is exceeded more than 100 per 5 seconds, it might possibly abusing.
	// so it must be limited.
	v, ok := h.LocalCache.Get(ip)
	if v != nil && v.(int64) > 100 {
		daily, total, err := h.Counter.GetHitOfDailyAndTotal(c.Request().Context(), id, time.Now())
		if err != nil {
			return err
		}
		return h.responseBadge(hctx, daily, total, title, titleBg,
			countBg, edgeType, icon, iconColor)
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
		daily, total, err := h.Counter.GetHitOfDailyAndTotal(c.Request().Context(), id, time.Now())
		if err != nil {
			return err
		}
		return h.responseBadge(hctx, daily, total, title, titleBg,
			countBg, edgeType, icon, iconColor)
	}
	h.LocalCache.Set(temporaryId, int64(1), 1*time.Second)

	daily, err := h.Counter.IncreaseHitOfDaily(c.Request().Context(), id, time.Now())
	if err != nil {
		return err
	}

	total, err := h.Counter.IncreaseHitOfTotal(c.Request().Context(), id)
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
		internal.SentryError(err)
	}

	// Broadcast message to users to which connected
	h.WebSocketBreaker.BroadCast(&WebSocketMessage{
		Payload: []byte(fmt.Sprintf("[%s] %s", internal.TimeToString(time.Now()), id))},
	)
	return h.responseBadge(hctx, daily, total, title, titleBg,
		countBg, edgeType, icon, iconColor)
}

// KeepCount is API, which it is not to increase page count.
func (h *Handler) KeepCount(c echo.Context) error {
	hctx := c.(*handler.HitCounterContext)
	if hctx.Get("ckid") == nil || hctx.Get("host") == nil || hctx.Get("path") == nil ||
		hctx.Get("title") == nil || hctx.Get("title_bg") == nil || hctx.Get("count_bg") == nil ||
		hctx.Get("edge_flat") == nil || hctx.Get("icon") == nil || hctx.Get("icon_color") == nil {
		return fmt.Errorf("[err] KeepCount empty params")
	}
	host := hctx.Get("host").(string)
	path := hctx.Get("path").(string)
	cookie := hctx.Get("ckid").(string)
	title := hctx.Get("title").(string)
	titleBg := hctx.Get("title_bg").(string)
	countBg := hctx.Get("count_bg").(string)
	edgeType := hctx.Get("edge_flat").(bool)
	icon := hctx.Get("icon").(string)
	iconColor := hctx.Get("icon_color").(string)

	_ = cookie
	id := fmt.Sprintf(countIdFormat, host, path)
	daily, total, err := h.Counter.GetHitOfDailyAndTotal(c.Request().Context(), id, time.Now())
	if err != nil {
		return err
	}

	return h.responseBadge(hctx, daily, total, title, titleBg,
		countBg, edgeType, icon, iconColor)
}

// DailyHitsInRecently is API, which shows a graph related daily page count.
func (h *Handler) DailyHitsInRecently(c echo.Context) error {
	hctx := c.(*handler.HitCounterContext)
	if hctx.Get("ckid") == nil || hctx.Get("host") == nil || hctx.Get("path") == nil {
		return fmt.Errorf("[err] KeepCount empty params")
	}
	host := hctx.Get("host").(string)
	path := hctx.Get("path").(string)
	cookie := hctx.Get("ckid").(string)
	_ = cookie

	var dateRange []time.Time
	now := time.Now()
	prev := time.Now().Add(-180 * 24 * time.Hour)
	for now.Unix() >= prev.Unix() {
		dateRange = append(dateRange, prev)
		prev = prev.Add(24 * time.Hour)
	}

	id := fmt.Sprintf(countIdFormat, host, path)
	scores, err := h.Counter.GetHitOfDailyByRange(c.Request().Context(), id, dateRange)
	if err != nil {
		return err
	}

	var yValues []float64
	for _, score := range scores {
		if score == nil {
			yValues = append(yValues, 0)
		} else {
			yValues = append(yValues, float64(score.Value))
		}
	}
	graph := chart.Chart{
		Width:      650,
		Height:     300,
		Title:      fmt.Sprintf("%s", id),
		TitleStyle: chart.StyleShow(),
		XAxis: chart.XAxis{
			Name:  "date",
			Style: chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:  "count",
			Style: chart.StyleShow(),
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: dateRange,
				YValues: yValues,
			},
		},
	}

	buf := new(bytes.Buffer)
	hctx.Response().Header().Set("Content-Type", chart.ContentTypeSVG)
	graph.Render(chart.SVG, buf)
	return hctx.String(http.StatusOK, string(buf.Bytes()))
}

func (h *Handler) responseBadge(ctx *handler.HitCounterContext,
	daily, total *counter.Score, titleText, titleBgColor, countBgColor string, edgeFlat bool,
	icon, iconColor string) error {
	dailyCount := int64(0)
	totalCount := int64(0)
	if daily != nil {
		dailyCount = daily.Value
	}
	if total != nil {
		totalCount = total.Value
	}

	// create default title
	if strings.TrimSpace(titleText) == "" {
		titleText = "hits"
	}

	// create default title background color
	if strings.TrimSpace(titleBgColor) == "" {
		titleBgColor = "#555"
	}

	// create default count background color
	if strings.TrimSpace(countBgColor) == "" {
		countBgColor = "#79c83d"
	}

	// create count text
	countText := fmt.Sprintf(badgeFormat, dailyCount, totalCount)

	// create badge
	var svg []byte
	var err error
	badge := internal.GenerateBadge(titleText, titleBgColor, countText, countBgColor, edgeFlat)
	if _, ok := h.Icons[icon]; !ok {
		svg, err = h.Badge.RenderFlatBadge(badge)
		if err != nil {
			return fmt.Errorf("[err] responseBadge %w", err)
		}
	} else {
		svg, err = h.Badge.RenderIconBadge(badge, icon, iconColor)
		if err != nil {
			return fmt.Errorf("[err] responseBadge %w", err)
		}
	}

	ctx.Response().Header().Set("Content-Type", "image/svg+xml")
	ctx.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Response().Header().Set("Pragma", "no-cache")
	ctx.Response().Header().Set("Expires", "0")
	return ctx.String(http.StatusOK, string(svg))
}
