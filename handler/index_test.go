package handler

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Index(t *testing.T) {
	assert := assert.New(t)
	defer mockRedis.FlushAll()


	e := echo.New()
	h, err := NewHandler(mockRedis.Addr())
	assert.NoError(err)

	ctx := context.Background()
	_, err = h.Counter.IncreaseRankOfTotal(ctx, "github.com", "/gjbae1212/hit-counter/")
	assert.NoError(err)
	_, err = h.Counter.IncreaseRankOfTotal(ctx, "github.com", "/gjbae1212/helloworld")
	assert.NoError(err)
	_, err = h.Counter.IncreaseRankOfTotal(ctx, "github.com", "/gjbae1212/power/dfdsfhtp(s///sdfsdf)")
	assert.NoError(err)

	tests := map[string]struct {
		included []string
		excluded []string
	}{
		"sample": {
			included: []string{"github.com/gjbae1212/hit-counter", "github.com/gjbae1212/helloworld"},
			excluded: []string{"github.com/gjbae1212/power"},
		},
	}

	for _, t := range tests {
		r := httptest.NewRequest("GET", "http://localhost:8080", nil)
		w := httptest.NewRecorder()
		hctx := &HitCounterContext{Context: e.NewContext(r, w)}

		err = h.Index(hctx)
		assert.NoError(err)

		resp := w.Result()
		assert.Equal(http.StatusOK, resp.StatusCode)
		raw, err := ioutil.ReadAll(resp.Body)
		assert.NoError(err)
		body := string(raw)

		for _, match := range t.included {
			assert.True(strings.Contains(body, match))
		}

		for _, match := range t.excluded {
			assert.False(strings.Contains(body, match))
		}
	}

}
