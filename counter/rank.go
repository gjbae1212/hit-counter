package counter

import (
	"fmt"
	"time"

	"strconv"

	allan_util "github.com/gjbae1212/go-module/util"
	"github.com/pkg/errors"
)

var (
	rankDailyFormat = "rank:daily:%s:%s"
	rankTotalFormat = "rank:total:%s"
)

func (d *db) IncreaseRankOfDaily(group, id string, t time.Time) (*Score, error) {
	if group == "" || id == "" || t.IsZero() {
		return nil, fmt.Errorf("[err] IncreaseRankOfDaily empty param")
	}

	daily := allan_util.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(rankDailyFormat, daily, group)
	v, err := d.redis.DoWithTimeout(timeout, "ZINCRBY", key, 1, id)
	if err != nil {
		return nil, errors.Wrap(err, "[err] IncreaseRankOfDaily")
	}

	rt, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "[err] IncreaseRankOfDaily")
	}

	return &Score{Name: id, Value: rt}, nil
}

func (d *db) IncreaseRankOfTotal(group, id string) (*Score, error) {
	if group == "" || id == "" {
		return nil, fmt.Errorf("[err] IncreaseRankOfTotal empty param")
	}

	key := fmt.Sprintf(rankTotalFormat, group)
	v, err := d.redis.DoWithTimeout(timeout, "ZINCRBY", key, 1, id)
	if err != nil {
		return nil, errors.Wrap(err, "[err] IncreaseRankOfTotal")
	}

	rt, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "[err] IncreaseRankOfTotal")
	}

	return &Score{Name: id, Value: rt}, nil
}

func (d *db) GetRankDailyByLimit(group string, limit int, t time.Time) ([]*Score, error) {
	if group == "" || limit <= 0 {
		return nil, fmt.Errorf("[err] GetRankDailyByLimit invalid param")
	}

	daily := allan_util.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(rankDailyFormat, daily, group)

	scores, err := d.redis.DoWithTimeout(timeout, "ZREVRANGE", key, 0, limit-1, "WITHSCORES")
	if err != nil {
		return nil, errors.Wrap(err, "[err] GetRankDailyByLimit")
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
			return nil, errors.Wrap(err, "[err] GetRankDailyByLimit")
		}
		rt = append(rt, &Score{Name: name, Value: v})
	}
	return rt, nil
}

func (d *db) GetRankTotalByLimit(group string, limit int) ([]*Score, error) {
	if group == "" || limit <= 0 {
		return nil, fmt.Errorf("[err] GetRankTotalByLimit invalid param")
	}

	key := fmt.Sprintf(rankTotalFormat, group)
	scores, err := d.redis.DoWithTimeout(timeout, "ZREVRANGE", key, 0, limit-1, "WITHSCORES")
	if err != nil {
		return nil, errors.Wrap(err, "[err] GetRankTotalByLimit")
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
			return nil, errors.Wrap(err, "[err] GetRankTotalByLimit")
		}

		rt = append(rt, &Score{Name: name, Value: v})
	}
	return rt, nil
}
