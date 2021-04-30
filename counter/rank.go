package counter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/gjbae1212/hit-counter/internal"
)

var (
	rankDailyFormat = "rank:daily:%s:%s"
	rankTotalFormat = "rank:total:%s"
)

// IncreaseRankOfDaily increases daily rank score.
func (d *db) IncreaseRankOfDaily(ctx context.Context, group, id string, t time.Time) (*Score, error) {
	if group == "" || id == "" || t.IsZero() {
		return nil, fmt.Errorf("[err] IncreaseRankOfDaily %w", internal.ErrorEmptyParams)
	}

	pipe := d.redisClient.Pipeline()

	daily := internal.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(rankDailyFormat, daily, group)

	// expire 2 month.
	incrResult := pipe.ZIncrBy(ctx, key, 1, id)
	pipe.Expire(ctx, key, time.Hour*24*60)

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, fmt.Errorf("[err] IncreaseRankOfDaily %w", err)
	}

	incr, _ := incrResult.Result()
	return &Score{Name: id, Value: int64(incr)}, nil
}

// IncreaseRankOfTotal increases accumulate rank score.
func (d *db) IncreaseRankOfTotal(ctx context.Context, group, id string) (*Score, error) {
	if group == "" || id == "" {
		return nil, fmt.Errorf("[err] IncreaseRankOfTotal %w", internal.ErrorEmptyParams)
	}

	key := fmt.Sprintf(rankTotalFormat, group)
	v, err := d.redisClient.ZIncrBy(ctx, key, 1, id).Result()
	if err != nil {
		return nil, fmt.Errorf("[err] IncreaseRankOfTotal %w", err)
	}

	return &Score{Name: id, Value: int64(v)}, nil
}

// GetRankDailyByLimit returns daily rank scores by limit.
func (d *db) GetRankDailyByLimit(ctx context.Context, group string, limit int, t time.Time) ([]*Score, error) {
	if group == "" || limit <= 0 {
		return nil, fmt.Errorf("[err] GetRankDailyByLimit %w", internal.ErrorEmptyParams)
	}
	var ret []*Score

	daily := internal.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(rankDailyFormat, daily, group)

	scores, err := d.redisClient.ZRevRangeWithScores(ctx, key, 0, int64(limit-1)).Result()
	if err == redis.Nil {
		return ret, nil
	} else if err != nil {
		return nil, fmt.Errorf("[err] GetRankDailyByLimit %w", err)
	}

	for _, score := range scores {
		name := score.Member.(string)
		value := score.Score
		ret = append(ret, &Score{Name: name, Value: int64(value)})
	}

	return ret, nil
}

// GetRankTotalByLimit returns total ranks.
func (d *db) GetRankTotalByLimit(ctx context.Context, group string, limit int) ([]*Score, error) {
	if group == "" || limit <= 0 {
		return nil, fmt.Errorf("[err] GetRankTotalByLimit %w", internal.ErrorEmptyParams)
	}

	var ret []*Score

	key := fmt.Sprintf(rankTotalFormat, group)

	scores, err := d.redisClient.ZRevRangeWithScores(ctx, key, 0, int64(limit-1)).Result()
	if err == redis.Nil {
		return ret, nil
	} else if err != nil {
		return nil, fmt.Errorf("[err] GetRankTotalByLimit %w", err)
	}

	for _, score := range scores {
		name := score.Member.(string)
		value := score.Score
		ret = append(ret, &Score{Name: name, Value: int64(value)})
	}

	return ret, nil
}
