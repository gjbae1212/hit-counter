package counter

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gjbae1212/hit-counter/internal"
)

var (
	rankDailyFormat = "rank:daily:%s:%s"
	rankTotalFormat = "rank:total:%s"
)

// IncreaseRankOfDaily increases daily rank score.
func (d *db) IncreaseRankOfDaily(group, id string, t time.Time) (*Score, error) {
	if group == "" || id == "" || t.IsZero() {
		return nil, fmt.Errorf("[err] IncreaseRankOfDaily %w", internal.ErrorEmptyParams)
	}

	daily := internal.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(rankDailyFormat, daily, group)
	v, err := d.redis.DoWithTimeout(timeout, "ZINCRBY", key, 1, id)
	if err != nil {
		return nil, fmt.Errorf("[err] IncreaseRankOfDaily %w", err)
	}

	rt, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("[err] IncreaseRankOfDaily %w", err)
	}

	return &Score{Name: id, Value: rt}, nil
}

// IncreaseRankOfTotal increases accumulate rank score.
func (d *db) IncreaseRankOfTotal(group, id string) (*Score, error) {
	if group == "" || id == "" {
		return nil, fmt.Errorf("[err] IncreaseRankOfTotal %w", internal.ErrorEmptyParams)
	}

	key := fmt.Sprintf(rankTotalFormat, group)
	v, err := d.redis.DoWithTimeout(timeout, "ZINCRBY", key, 1, id)
	if err != nil {
		return nil, fmt.Errorf("[err] IncreaseRankOfTotal %w", err)
	}

	rt, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("[err] IncreaseRankOfTotal %w", err)
	}

	return &Score{Name: id, Value: rt}, nil
}

//  GetRankDailyByLimit returns daily rank scores by limit.
func (d *db) GetRankDailyByLimit(group string, limit int, t time.Time) ([]*Score, error) {
	if group == "" || limit <= 0 {
		return nil, fmt.Errorf("[err] GetRankDailyByLimit %w", internal.ErrorEmptyParams)
	}

	daily := internal.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(rankDailyFormat, daily, group)

	scores, err := d.redis.DoWithTimeout(timeout, "ZREVRANGE", key, 0, limit-1, "WITHSCORES")
	if err != nil {
		return nil, fmt.Errorf("[err] GetRankDailyByLimit %w", err)
	}

	var rt []*Score
	// empty
	if scores == nil {
		return rt, nil
	}

	list := scores.([]interface{})
	for i := 0; i < len(list); i += 2 {
		name := string(list[i].([]byte))
		value := string(list[i+1].([]byte))
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("[err] GetRankDailyByLimit %w", err)
		}
		rt = append(rt, &Score{Name: name, Value: v})
	}
	return rt, nil
}

// GetRankTotalByLimit returns total ranks.
func (d *db) GetRankTotalByLimit(group string, limit int) ([]*Score, error) {
	if group == "" || limit <= 0 {
		return nil, fmt.Errorf("[err] GetRankTotalByLimit %w", internal.ErrorEmptyParams)
	}

	key := fmt.Sprintf(rankTotalFormat, group)
	scores, err := d.redis.DoWithTimeout(timeout, "ZREVRANGE", key, 0, limit-1, "WITHSCORES")
	if err != nil {
		return nil, fmt.Errorf("[err] GetRankTotalByLimit %w", err)
	}

	var rt []*Score
	// empty
	if scores == nil {
		return rt, nil
	}

	list := scores.([]interface{})
	for i := 0; i < len(list); i += 2 {
		name := string(list[i].([]byte))
		value := string(list[i+1].([]byte))
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("[err] GetRankTotalByLimit %w", err)
		}

		rt = append(rt, &Score{Name: name, Value: v})
	}
	return rt, nil
}
