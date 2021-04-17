package api_handler

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gjbae1212/hit-counter/handler"
	"github.com/gjbae1212/hit-counter/internal"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_KeepCount(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	h, err := handler.NewHandler(mockRedis.Addr())
	assert.NoError(err)

	api, err := NewHandler(h)
	assert.NoError(err)
	e := echo.New()

	// err
	errR := httptest.NewRequest("GET", "http://localhost:8080", nil)
	errW := httptest.NewRecorder()
	errCtx := &handler.HitCounterContext{Context: e.NewContext(errR, errW)}

	// default
	defaultR := httptest.NewRequest("GET", "http://localhost:8080", nil)
	defaultW := httptest.NewRecorder()
	defaultCtx := &handler.HitCounterContext{Context: e.NewContext(defaultR, defaultW)}
	defaultCtx.Set("ckid", "test")
	defaultCtx.Set("host", "github.com")
	defaultCtx.Set("path", "gjbae1212/hit-counter")
	defaultCtx.Set("title", " ")
	defaultCtx.Set("title_bg", " ")
	defaultCtx.Set("count_bg", " ")
	defaultCtx.Set("edge_flat", true)
	defaultCtx.Set("icon", "")
	defaultCtx.Set("icon_color", "")

	defaultOutput, err := h.Badge.RenderFlatBadge(internal.GenerateBadge("hits",
		"#555", fmt.Sprintf(badgeFormat, 0, 0), "#79c83d", true))
	assert.NoError(err)

	// title
	titleR := httptest.NewRequest("GET", "http://localhost:8080", nil)
	titleW := httptest.NewRecorder()
	titleCtx := &handler.HitCounterContext{Context: e.NewContext(titleR, titleW)}
	titleCtx.Set("ckid", "test")
	titleCtx.Set("host", "github.com")
	titleCtx.Set("path", "gjbae1212/hit-counter")
	titleCtx.Set("title", " hello ")
	titleCtx.Set("title_bg", "")
	titleCtx.Set("count_bg", "")
	titleCtx.Set("edge_flat", true)
	titleCtx.Set("icon", "")
	titleCtx.Set("icon_color", "")

	titleOutput, err := h.Badge.RenderFlatBadge(internal.GenerateBadge(" hello ",
		"#555", fmt.Sprintf(badgeFormat, 0, 0), "#79c83d", true))
	assert.NoError(err)

	// bg-color
	bgColorR := httptest.NewRequest("GET", "http://localhost:8080", nil)
	bgColorW := httptest.NewRecorder()
	bgColorCtx := &handler.HitCounterContext{Context: e.NewContext(bgColorR, bgColorW)}
	bgColorCtx.Set("ckid", "test")
	bgColorCtx.Set("host", "github.com")
	bgColorCtx.Set("path", "gjbae1212/hit-counter")
	bgColorCtx.Set("title", "")
	bgColorCtx.Set("title_bg", "#111")
	bgColorCtx.Set("count_bg", "#222")
	bgColorCtx.Set("edge_flat", true)
	bgColorCtx.Set("icon", "")
	bgColorCtx.Set("icon_color", "")

	bgColorOutput, err := h.Badge.RenderFlatBadge(internal.GenerateBadge("hits",
		"#111", fmt.Sprintf(badgeFormat, 0, 0), "#222", true))
	assert.NoError(err)

	// edge
	edgeR := httptest.NewRequest("GET", "http://localhost:8080", nil)
	edgeW := httptest.NewRecorder()
	edgeCtx := &handler.HitCounterContext{Context: e.NewContext(edgeR, edgeW)}
	edgeCtx.Set("ckid", "test")
	edgeCtx.Set("host", "github.com")
	edgeCtx.Set("path", "gjbae1212/hit-counter")
	edgeCtx.Set("title", "")
	edgeCtx.Set("title_bg", "")
	edgeCtx.Set("count_bg", "")
	edgeCtx.Set("edge_flat", false)
	edgeCtx.Set("icon", "")
	edgeCtx.Set("icon_color", "")

	edgeCtxOutput, err := h.Badge.RenderFlatBadge(internal.GenerateBadge("hits",
		"#555", fmt.Sprintf(badgeFormat, 0, 0), "#79c83d", false))
	assert.NoError(err)

	// icon
	iconR := httptest.NewRequest("GET", "http://localhost:8080", nil)
	iconW := httptest.NewRecorder()
	iconCtx := &handler.HitCounterContext{Context: e.NewContext(iconR, iconW)}
	iconCtx.Set("ckid", "test")
	iconCtx.Set("host", "github.com")
	iconCtx.Set("path", "gjbae1212/hit-counter")
	iconCtx.Set("title", "")
	iconCtx.Set("title_bg", "")
	iconCtx.Set("count_bg", "")
	iconCtx.Set("edge_flat", false)
	iconCtx.Set("icon", "a-frame.svg")
	iconCtx.Set("icon_color", "")

	iconCtxOutput, err := h.Badge.RenderIconBadge(internal.GenerateBadge("hits",
		"#555", fmt.Sprintf(badgeFormat, 0, 0), "#79c83d", false), "a-frame.svg", "")
	assert.NoError(err)

	// icon with color
	iconWithColorR := httptest.NewRequest("GET", "http://localhost:8080", nil)
	iconWithColorW := httptest.NewRecorder()
	iconWithColorCtx := &handler.HitCounterContext{Context: e.NewContext(iconWithColorR, iconWithColorW)}
	iconWithColorCtx.Set("ckid", "test")
	iconWithColorCtx.Set("host", "github.com")
	iconWithColorCtx.Set("path", "gjbae1212/hit-counter")
	iconWithColorCtx.Set("title", "")
	iconWithColorCtx.Set("title_bg", "")
	iconWithColorCtx.Set("count_bg", "")
	iconWithColorCtx.Set("edge_flat", false)
	iconWithColorCtx.Set("icon", "a-frame.svg")
	iconWithColorCtx.Set("icon_color", "#aaaaaa")

	iconWithColorCtxOutput, err := h.Badge.RenderIconBadge(internal.GenerateBadge("hits",
		"#555", fmt.Sprintf(badgeFormat, 0, 0), "#79c83d", false), "a-frame.svg", "#aaaaaa")
	assert.NoError(err)

	tests := map[string]struct {
		input  *handler.HitCounterContext
		w      *httptest.ResponseRecorder
		output string
		isErr  bool
	}{
		"err": {
			input: errCtx,
			isErr: true,
		},
		"default": {
			input:  defaultCtx,
			w:      defaultW,
			output: string(defaultOutput),
		},
		"title": {
			input:  titleCtx,
			w:      titleW,
			output: string(titleOutput),
		},
		"bg-color": {
			input:  bgColorCtx,
			w:      bgColorW,
			output: string(bgColorOutput),
		},
		"edge": {
			input:  edgeCtx,
			w:      edgeW,
			output: string(edgeCtxOutput),
		},
		"icon": {
			input:  iconCtx,
			w:      iconW,
			output: string(iconCtxOutput),
		},
		"icon-with-color": {
			input:  iconWithColorCtx,
			w:      iconWithColorW,
			output: string(iconWithColorCtxOutput),
		},
	}

	for _, t := range tests {
		err := api.KeepCount(t.input)
		assert.Equal(t.isErr, err != nil)
		if err == nil {
			assert.Equal(200, t.w.Code)
			assert.Equal(t.output, t.w.Body.String())
		}
	}
}

func TestHandler_IncrCount(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()

	h, err := handler.NewHandler(mockRedis.Addr())
	assert.NoError(err)

	api, err := NewHandler(h)
	assert.NoError(err)
	e := echo.New()

	// err
	errR := httptest.NewRequest("GET", "http://localhost:8080", nil)
	errW := httptest.NewRecorder()
	errCtx := &handler.HitCounterContext{Context: e.NewContext(errR, errW)}

	// default
	defaultR := httptest.NewRequest("GET", "http://localhost:8080", nil)
	defaultW := httptest.NewRecorder()
	defaultCtx := &handler.HitCounterContext{Context: e.NewContext(defaultR, defaultW)}
	defaultCtx.Set("ckid", "test")
	defaultCtx.Set("host", "github.com")
	defaultCtx.Set("path", "gjbae1212/hit-counter-default")
	defaultCtx.Set("title", " ")
	defaultCtx.Set("title_bg", " ")
	defaultCtx.Set("count_bg", " ")
	defaultCtx.Set("edge_flat", true)
	defaultCtx.Set("icon", "")
	defaultCtx.Set("icon_color", "")

	defaultOutput, err := h.Badge.RenderFlatBadge(internal.GenerateBadge("hits",
		"#555", fmt.Sprintf(badgeFormat, 1, 1), "#79c83d", true))
	assert.NoError(err)

	// title
	titleR := httptest.NewRequest("GET", "http://localhost:8080", nil)
	titleW := httptest.NewRecorder()
	titleCtx := &handler.HitCounterContext{Context: e.NewContext(titleR, titleW)}
	titleCtx.Set("ckid", "test")
	titleCtx.Set("host", "github.com")
	titleCtx.Set("path", "gjbae1212/hit-counter-title")
	titleCtx.Set("title", " hello ")
	titleCtx.Set("title_bg", "")
	titleCtx.Set("count_bg", "")
	titleCtx.Set("edge_flat", true)
	titleCtx.Set("icon", "")
	titleCtx.Set("icon_color", "")

	titleOutput, err := h.Badge.RenderFlatBadge(internal.GenerateBadge(" hello ",
		"#555", fmt.Sprintf(badgeFormat, 1, 1), "#79c83d", true))
	assert.NoError(err)

	// bg-color
	bgColorR := httptest.NewRequest("GET", "http://localhost:8080", nil)
	bgColorW := httptest.NewRecorder()
	bgColorCtx := &handler.HitCounterContext{Context: e.NewContext(bgColorR, bgColorW)}
	bgColorCtx.Set("ckid", "test")
	bgColorCtx.Set("host", "github.com")
	bgColorCtx.Set("path", "gjbae1212/hit-counter-bg-color")
	bgColorCtx.Set("title", "")
	bgColorCtx.Set("title_bg", "#111")
	bgColorCtx.Set("count_bg", "#222")
	bgColorCtx.Set("edge_flat", true)
	bgColorCtx.Set("icon", "")
	bgColorCtx.Set("icon_color", "")

	bgColorOutput, err := h.Badge.RenderFlatBadge(internal.GenerateBadge("hits",
		"#111", fmt.Sprintf(badgeFormat, 1, 1), "#222", true))
	assert.NoError(err)

	// edge
	edgeR := httptest.NewRequest("GET", "http://localhost:8080", nil)
	edgeW := httptest.NewRecorder()
	edgeCtx := &handler.HitCounterContext{Context: e.NewContext(edgeR, edgeW)}
	edgeCtx.Set("ckid", "test")
	edgeCtx.Set("host", "github.com")
	edgeCtx.Set("path", "gjbae1212/hit-counter-edge")
	edgeCtx.Set("title", "")
	edgeCtx.Set("title_bg", "")
	edgeCtx.Set("count_bg", "")
	edgeCtx.Set("edge_flat", false)
	edgeCtx.Set("icon", "")
	edgeCtx.Set("icon_color", "")

	edgeCtxOutput, err := h.Badge.RenderFlatBadge(internal.GenerateBadge("hits",
		"#555", fmt.Sprintf(badgeFormat, 1, 1), "#79c83d", false))
	assert.NoError(err)

	tests := map[string]struct {
		input  *handler.HitCounterContext
		w      *httptest.ResponseRecorder
		output string
		isErr  bool
	}{
		"err": {
			input: errCtx,
			isErr: true,
		},
		"default": {
			input:  defaultCtx,
			w:      defaultW,
			output: string(defaultOutput),
		},
		"title": {
			input:  titleCtx,
			w:      titleW,
			output: string(titleOutput),
		},
		"bg-color": {
			input:  bgColorCtx,
			w:      bgColorW,
			output: string(bgColorOutput),
		},
		"edge": {
			input:  edgeCtx,
			w:      edgeW,
			output: string(edgeCtxOutput),
		},
	}

	for k, t := range tests {
		t.input.Request().Header.Set("User-Agent", k)
		err := api.IncrCount(t.input)
		assert.Equal(t.isErr, err != nil)
		if err == nil {
			assert.Equal(200, t.w.Code)
			assert.Equal(t.output, t.w.Body.String())
		}
	}

	for i := 0; i < 10; i++ {
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		r.Header.Set("User-Agent", fmt.Sprintf("%d", i))
		w := httptest.NewRecorder()
		hctx := &handler.HitCounterContext{Context: e.NewContext(r, w)}
		hctx.Set("ckid", "test")
		hctx.Set("host", "github.com")
		hctx.Set("path", "gjbae1212/hit-counter")
		hctx.Set("title", "")
		hctx.Set("title_bg", "")
		hctx.Set("count_bg", "")
		hctx.Set("edge_flat", false)
		hctx.Set("icon", "")
		hctx.Set("icon_color", "")

		err = api.IncrCount(hctx)
		assert.NoError(err)
		assert.Equal(200, w.Code)
	}

	time.Sleep(3 * time.Second)
	scores, err := api.Counter.GetRankDailyByLimit(context.Background(), "github.com", 10, time.Now())
	assert.NoError(err)
	assert.Len(scores, 5)
	assert.Equal(int64(10), scores[0].Value)
	scores, err = api.Counter.GetRankTotalByLimit(context.Background(), "github.com", 10)
	assert.NoError(err)
	assert.Len(scores, 5)
	assert.Equal(int64(10), scores[0].Value)
	scores, err = api.Counter.GetRankTotalByLimit(context.Background(), "domain", 10)
	assert.NoError(err)
	assert.Len(scores, 1)
	assert.Equal(int64(14), scores[0].Value)
}
