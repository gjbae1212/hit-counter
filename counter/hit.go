package counter

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/gjbae1212/hit-counter/internal"
)

var (
	hitDailyFormat = "hit:daily:%s:%s"
	hitTotalFormat = "hit:total:%s"
)

// IncreaseHitOfDaily increases daily count.
func (d *db) IncreaseHitOfDaily(ctx context.Context, id string, t time.Time) (*Score, error) {
	if id == "" || t.IsZero() {
		return nil, fmt.Errorf("[err] IncreaseHitOfDaily  %w", internal.ErrorEmptyParams)
	}

	pipe := d.redisClient.Pipeline()
	daily := internal.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(hitDailyFormat, daily, id)

	// expire 2 month.
	incrResult := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Hour*24*60)

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, fmt.Errorf("[err] IncreaseHitOfDaily %w", err)
	}

	incr, _ := incrResult.Result()
	return &Score{Name: key, Value: incr}, nil
}

// IncreaseHitOfTotal increases accumulate count.
func (d *db) IncreaseHitOfTotal(ctx context.Context, id string) (*Score, error) {
	if id == "" {
		return nil, fmt.Errorf("[err] IncreaseHitOfTotal %w", internal.ErrorEmptyParams)
	}

	key := fmt.Sprintf(hitTotalFormat, id)
	v, err := d.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("[err] IncreaseHitOfTotal %w", err)
	}
	return &Score{Name: key, Value: v}, nil
}

// GetHitOfDaily returns daily score.
func (d *db) GetHitOfDaily(ctx context.Context, id string, t time.Time) (*Score, error) {
	if id == "" || t.IsZero() {
		return nil, fmt.Errorf("[err] GetHitOfDaily empty param")
	}

	daily := internal.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(hitDailyFormat, daily, id)

	v, err := d.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("[err] GetHitOfDaily %w", err)
	}

	rt, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("[err] GetHitOfDaily %w", err)
	}

	return &Score{Name: key, Value: rt}, nil
}

// GetHitOfTotal returns  accumulate score.
func (d *db) GetHitOfTotal(ctx context.Context, id string) (*Score, error) {
	if id == "" {
		return nil, fmt.Errorf("[err] GetHitOfTotal empty param")
	}

	key := fmt.Sprintf(hitTotalFormat, id)
	v, err := d.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("[err] GetHitOfTotal %w", err)
	}

	rt, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("[err] GetHitOfTotal %w", err)
	}

	return &Score{Name: key, Value: rt}, nil
}

// GetHitOfDailyAndTotal returns daily score and  accumulate score.
func (d *db) GetHitOfDailyAndTotal(ctx context.Context, id string, t time.Time) (daily *Score, total *Score, retErr error) {
	if id == "" || t.IsZero() {
		retErr = fmt.Errorf("[err] GetHitOfDailyAndTotal %w", internal.ErrorEmptyParams)
		return
	}

	key1 := fmt.Sprintf(hitDailyFormat, internal.TimeToDailyStringFormat(t), id)
	key2 := fmt.Sprintf(hitTotalFormat, id)

	v, err := d.redisClient.MGet(ctx, key1, key2).Result()
	if err == redis.Nil {
		return
	} else if err != nil {
		retErr = fmt.Errorf("[err] GetHitOfDailyAndTotal %w", err)
		return
	}

	if v[0] != nil {
		dailyValue, err := strconv.ParseInt(v[0].(string), 10, 64)
		if err != nil {
			retErr = fmt.Errorf("[err] GetHitOfDailyAndTotal %w", err)
			return
		}
		daily = &Score{Name: key1, Value: dailyValue}
	}

	if v[1] != nil {
		totalValue, err := strconv.ParseInt(v[1].(string), 10, 64)
		if err != nil {
			retErr = fmt.Errorf("[err] GetHitOfDailyAndTotal %w", err)
			return
		}
		total = &Score{Name: key2, Value: totalValue}
	}
	return
}

// GetHitOfDailyByRange returns daily scores with range.
func (d *db) GetHitOfDailyByRange(ctx context.Context, id string, timeRange []time.Time) (scores []*Score, retErr error) {
	if id == "" || len(timeRange) == 0 {
		retErr = fmt.Errorf("[err] GetHitOfDailyByRange %w", internal.ErrorEmptyParams)
		return
	}

	var keys []string
	for _, t := range timeRange {
		keys = append(keys, fmt.Sprintf(hitDailyFormat, internal.TimeToDailyStringFormat(t), id))
	}

	v, err := d.redisClient.MGet(ctx, keys...).Result()
	if err == redis.Nil {
		return
	} else if err != nil {
		retErr = fmt.Errorf("[err] GetHitOfDailyByRange %w", err)
	}

	for i, key := range keys {
		if v[i] != nil {
			dailyValue, err := strconv.ParseInt(v[i].(string), 10, 64)
			if err != nil {
				err = fmt.Errorf("[err] GetHitOfDailyByRange %w", err)
				return
			}
			scores = append(scores, &Score{Name: key, Value: dailyValue})
		} else {
			scores = append(scores, nil)
		}
	}

	return
}
