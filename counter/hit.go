package counter

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gjbae1212/hit-counter/internal"
)

var (
	hitDailyFormat = "hit:daily:%s:%s"
	hitTotalFormat = "hit:total:%s"
)

// IncreaseHitOfDaily increases daily count.
func (d *db) IncreaseHitOfDaily(id string, t time.Time) (*Score, error) {
	if id == "" || t.IsZero() {
		return nil, fmt.Errorf("[err] IncreaseHitOfDaily  %w", internal.ErrorEmptyParams)
	}

	daily := internal.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(hitDailyFormat, daily, id)
	v, err := d.redis.DoWithTimeout(timeout, "INCR", key)
	if err != nil {
		return nil, fmt.Errorf("[err] IncreaseHitOfDaily %w", err)
	}

	return &Score{Name: key, Value: v.(int64)}, nil
}

// IncreaseHitOfTotal increases accumulate count.
func (d *db) IncreaseHitOfTotal(id string) (*Score, error) {
	if id == "" {
		return nil, fmt.Errorf("[err] IncreaseHitOfTotal %w", internal.ErrorEmptyParams)
	}

	key := fmt.Sprintf(hitTotalFormat, id)
	v, err := d.redis.DoWithTimeout(timeout, "INCR", key)
	if err != nil {
		return nil, fmt.Errorf("[err] IncreaseHitOfTotal %w", err)
	}
	return &Score{Name: key, Value: v.(int64)}, nil
}

// GetHitOfDaily returns daily score.
func (d *db) GetHitOfDaily(id string, t time.Time) (*Score, error) {
	if id == "" || t.IsZero() {
		return nil, fmt.Errorf("[err] GetHitOfDaily empty param")
	}

	daily := internal.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(hitDailyFormat, daily, id)

	v, err := d.redis.DoWithTimeout(timeout, "GET", key)
	if err != nil {
		return nil, fmt.Errorf("[err] GetHitOfDaily %w", err)
	}

	// if v is empty
	if v == nil {
		return nil, nil
	}

	rt, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("[err] GetHitOfDaily %w", err)
	}

	return &Score{Name: key, Value: rt}, nil
}

// GetHitOfDaily returns  accumulate score.
func (d *db) GetHitOfTotal(id string) (*Score, error) {
	if id == "" {
		return nil, fmt.Errorf("[err] GetHitOfTotal empty param")
	}

	key := fmt.Sprintf(hitTotalFormat, id)
	v, err := d.redis.DoWithTimeout(timeout, "GET", key)
	if err != nil {
		return nil, fmt.Errorf("[err] GetHitOfTotal %w", err)
	}

	// if v is empty
	if v == nil {
		return nil, nil
	}

	rt, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("[err] GetHitOfTotal %w", err)
	}

	return &Score{Name: key, Value: rt}, nil
}

// GetHitOfDailyAndTotal returns daily score and  accumulate score.
func (d *db) GetHitOfDailyAndTotal(id string, t time.Time) (daily *Score, total *Score, err error) {
	if id == "" || t.IsZero() {
		err = fmt.Errorf("[err] GetHitOfDailyAndTotal %w", internal.ErrorEmptyParams)
		return
	}

	key1 := fmt.Sprintf(hitDailyFormat, internal.TimeToDailyStringFormat(t), id)
	key2 := fmt.Sprintf(hitTotalFormat, id)

	v, suberr := d.redis.DoWithTimeout(timeout, "MGET", key1, key2)
	if suberr != nil {
		err = fmt.Errorf("[err] GetHitOfDailyAndTotal %w", suberr)
		return
	}

	if v.([]interface{})[0] != nil {
		dailyValue, suberr := strconv.ParseInt(string(v.([]interface{})[0].([]byte)), 10, 64)
		if suberr != nil {
			err = fmt.Errorf("[err] GetHitOfDailyAndTotal %w", suberr)
			return
		}
		daily = &Score{Name: key1, Value: dailyValue}
	}

	if v.([]interface{})[1] != nil {
		totalValue, suberr := strconv.ParseInt(string(v.([]interface{})[1].([]byte)), 10, 64)
		if suberr != nil {
			err = fmt.Errorf("[err] GetHitOfDailyAndTotal %w", suberr)
			return
		}
		total = &Score{Name: key2, Value: totalValue}
	}
	return
}

// GetHitOfDailyByRange returns daily scores with range.
func (d *db) GetHitOfDailyByRange(id string, timeRange []time.Time) (scores []*Score, err error) {
	if id == "" || len(timeRange) == 0 {
		err = fmt.Errorf("[err] GetHitOfDailyByRange %w", internal.ErrorEmptyParams)
		return
	}

	var keys []interface{}
	for _, t := range timeRange {
		keys = append(keys, fmt.Sprintf(hitDailyFormat, internal.TimeToDailyStringFormat(t), id))
	}

	v, suberr := d.redis.DoWithTimeout(timeout, "MGET", keys...)
	if suberr != nil {
		err = fmt.Errorf("[err] GetHitOfDailyByRange %w", suberr)
		return
	}

	for i, key := range keys {
		if v.([]interface{})[i] != nil {
			dailyValue, suberr := strconv.ParseInt(string(v.([]interface{})[i].([]byte)), 10, 64)
			if suberr != nil {
				err = fmt.Errorf("[err] GetHitOfDailyByRange %w", suberr)
				return
			}
			scores = append(scores, &Score{Name: key.(string), Value: dailyValue})
		} else {
			scores = append(scores, nil)
		}
	}
	return
}
