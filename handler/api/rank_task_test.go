package api_handler

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/davecgh/go-spew/spew"
	"github.com/gjbae1212/hit-counter/counter"
	"github.com/gjbae1212/hit-counter/handler"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestRankTask_Process(t *testing.T) {
	assert := assert.New(t)

	s, err := miniredis.Run()
	assert.NoError(err)
	defer s.Close()

	h, err := handler.NewHandler([]string{s.Addr()})
	assert.NoError(err)

	api, err := NewHandler(h)
	assert.NoError(err)

	tests := map[string]struct {
		input *RankTask
		wants []*counter.Score
	}{
		"not_github": {input: &RankTask{
			Counter:   h.Counter,
			Domain:    "allan.com",
			Path:      "aa/bb",
			CreatedAt: time.Now(),
		}, wants: []*counter.Score{
			nil,
			nil,
			&counter.Score{
				Name:  "allan.com",
				Value: 1,
			},
			&counter.Score{
				Name:  "allan.com",
				Value: 1,
			},
		}},
		"github-1": {input: &RankTask{
			Counter:   h.Counter,
			Domain:    "github.com",
			Path:      "gjbae1212/test",
			CreatedAt: time.Now(),
		}, wants: []*counter.Score{
			&counter.Score{
				Name:  "gjbae1212/test",
				Value: 1,
			},
			&counter.Score{
				Name:  "gjbae1212/test",
				Value: 1,
			},
			nil,
			nil,
		}},
		"github-2": {input: &RankTask{
			Counter:   h.Counter,
			Domain:    "github.com",
			Path:      "gjbae1212/hoho",
			CreatedAt: time.Now(),
		}, wants: []*counter.Score{
			&counter.Score{
				Name:  "gjbae1212/hoho",
				Value: 1,
			},
			&counter.Score{
				Name:  "gjbae1212/hoho",
				Value: 1,
			},
			nil,
			nil,
		}},
	}

	ctx := context.Background()

	for k, t := range tests {
		switch k {
		case "not_github":
			err = api.AsyncTask.AddTask(ctx, t.input)
			assert.NoError(err)
			time.Sleep(1 * time.Second)

			scores, err := api.Counter.GetRankDailyByLimit(t.input.Domain, 10, time.Now())
			assert.NoError(err)
			assert.Len(scores, 0)

			scores, err = api.Counter.GetRankTotalByLimit(t.input.Domain, 10)
			assert.NoError(err)
			assert.Len(scores, 0)

			scores, err = api.Counter.GetRankDailyByLimit(domainGroup, 10, time.Now())
			assert.NoError(err)
			assert.Len(scores, 1)
			assert.True(cmp.Equal(t.wants[2], scores[0]))

			scores, err = api.Counter.GetRankTotalByLimit(domainGroup, 10)
			assert.NoError(err)
			assert.Len(scores, 1)
			assert.True(cmp.Equal(t.wants[3], scores[0]))
		case "github-1", "github-2":
			err = api.AsyncTask.AddTask(ctx, t.input)
			assert.NoError(err)
			time.Sleep(1 * time.Second)

			scores, err := api.Counter.GetRankDailyByLimit(githubGroup, 10, time.Now())
			assert.NoError(err)

			if len(scores) == 1 {
				assert.True(cmp.Equal(t.wants[0], scores[0]))
				spew.Dump(scores)
				scores, err = api.Counter.GetRankTotalByLimit(githubGroup, 10)
				assert.NoError(err)
				assert.Len(scores, 1)
				assert.True(cmp.Equal(t.wants[1], scores[0]))
				spew.Dump(scores)
			} else {
				assert.Len(scores, 2)
			}
		}
	}

	// [TEST] github.com domain, profile
	scores, err := api.Counter.GetRankDailyByLimit(domainGroup, 10, time.Now())
	assert.NoError(err)
	assert.Len(scores, 2)
	assert.True(cmp.Equal(&counter.Score{Name: githubGroup, Value: 2}, scores[0]))

	scores, err = api.Counter.GetRankTotalByLimit(domainGroup, 10)
	assert.NoError(err)
	assert.Len(scores, 2)
	assert.True(cmp.Equal(&counter.Score{Name: githubGroup, Value: 2}, scores[0]))

	scores, err = api.Counter.GetRankDailyByLimit(githubProfileSumGroup, 10, time.Now())
	assert.NoError(err)
	assert.Len(scores, 1)
	assert.True(cmp.Equal(&counter.Score{Name: "gjbae1212", Value: 2}, scores[0]))


}
