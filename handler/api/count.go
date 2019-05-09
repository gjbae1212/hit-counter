package api_handler

import (
	"net/http"

	"time"

	"fmt"

	"github.com/gjbae1212/go-badge"
	"github.com/gjbae1212/hit-counter/handler"
	"github.com/labstack/echo/v4"
)

var (
	badgeFormat = " %d / %d "
)

// TODO: 어뷰징 체크 (freecache 이용 xxhash로 저장)
// TODO: cache 옵션 끄기, content type 정리
func (h *Handler) IncrCount(c echo.Context) error {
	return c.String(http.StatusOK, "health check!")
}

func (h *Handler) KeepCount(c echo.Context) error {
	hctx := c.(*handler.HitCounterContext)
	id := hctx.Get("id").(string)

	daily, total, err := h.Counter.GetHitAll(id, time.Now())
	if err != nil {
		return err
	}

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

	c.Response().Header().Set("Content-Type", "image/svg+xml")
	c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("Expires", "0")
	return c.String(http.StatusOK, string(badge))
}
