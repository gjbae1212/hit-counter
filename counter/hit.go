package counter

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	allan_util "github.com/gjbae1212/go-module/util"
)

var (
	hitDailyFormat = "hit:daily:%s:%s"
	hitTotalFormat = "hit:total:%s"
)

func (d *db) IncreaseHitOfDaily(id string, t time.Time) (*Score, error) {
	if id == "" || t.IsZero() {
		return nil, fmt.Errorf("[err] IncreaseHitOfDaily empty param")
	}

	daily := allan_util.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(hitDailyFormat, daily, id)
	v, err := d.redis.DoWithTimeout(timeout, "INCR", key)
	if err != nil {
		return nil, errors.Wrap(err, "[err] IncreaseHitOfDaily")
	}

	return &Score{Name: key, Value: v.(int64)}, nil
}

func (d *db) IncreaseHitOfTotal(id string) (*Score, error) {
	if id == "" {
		return nil, fmt.Errorf("[err] IncreaseHitOfTotal empty param")
	}

	key := fmt.Sprintf(hitTotalFormat, id)
	v, err := d.redis.DoWithTimeout(timeout, "INCR", key)
	if err != nil {
		return nil, errors.Wrap(err, "[err] IncreaseHitOfTotal")
	}
	return &Score{Name: key, Value: v.(int64)}, nil
}

func (d *db) GetHitOfDaily(id string, t time.Time) (*Score, error) {
	if id == "" || t.IsZero() {
		return nil, fmt.Errorf("[err] GetHitOfDaily empty param")
	}

	daily := allan_util.TimeToDailyStringFormat(t)
	key := fmt.Sprintf(hitDailyFormat, daily, id)

	v, err := d.redis.DoWithTimeout(timeout, "GET", key)
	if err != nil {
		return nil, errors.Wrap(err, "[err] GetHitOfDaily")
	}

	// empty
	if v == nil {
		return nil, nil
	}

	rt, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "[err] GetHitOfDaily")
	}

	return &Score{Name: key, Value: rt}, nil
}

func (d *db) GetHitOfTotal(id string) (*Score, error) {
	if id == "" {
		return nil, fmt.Errorf("[err] GetHitOfTotal empty param")
	}

	key := fmt.Sprintf(hitTotalFormat, id)
	v, err := d.redis.DoWithTimeout(timeout, "GET", key)
	if err != nil {
		return nil, errors.Wrap(err, "[err] GetHitOfTotal")
	}

	// empty
	if v == nil {
		return nil, nil
	}

	rt, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "[err] GetHitOfTotal")
	}

	return &Score{Name: key, Value: rt}, nil
}

func (d *db) GetHitOfDailyAndTotal(id string, t time.Time) (daily *Score, total *Score, err error) {
	if id == "" || t.IsZero() {
		err = fmt.Errorf("[err] GetHitAll empty param")
		return
	}

	key1 := fmt.Sprintf(hitDailyFormat, allan_util.TimeToDailyStringFormat(t), id)
	key2 := fmt.Sprintf(hitTotalFormat, id)

	v, suberr := d.redis.DoWithTimeout(timeout, "MGET", key1, key2)
	if suberr != nil {
		err = errors.Wrap(suberr, "[err] GetHitAll")
		return
	}

	if v.([]interface{})[0] != nil {
		dailyValue, suberr := strconv.ParseInt(string(v.([]interface{})[0].([]byte)), 10, 64)
		if suberr != nil {
			err = errors.Wrap(suberr, "[err] GetHitAll")
			return
		}
		daily = &Score{Name: key1, Value: dailyValue}
	}

	if v.([]interface{})[1] != nil {
		totalValue, suberr := strconv.ParseInt(string(v.([]interface{})[1].([]byte)), 10, 64)
		if suberr != nil {
			err = errors.Wrap(suberr, "[err] GetHitAll")
			return
		}
		total = &Score{Name: key2, Value: totalValue}
	}
	return
}

func (d *db) GetHitOfDailyByRange(id string, timeRange []time.Time) (scores []*Score, err error) {
	if id == "" || len(timeRange) == 0 {
		err = fmt.Errorf("[err] GetHitOfDailyByRange")
		return
	}

	var keys []interface{}
	for _, t := range timeRange {
		keys = append(keys, fmt.Sprintf(hitDailyFormat, allan_util.TimeToDailyStringFormat(t), id))
	}

	v, suberr := d.redis.DoWithTimeout(timeout, "MGET", keys...)
	if suberr != nil {
		err = errors.Wrap(suberr, "[err] GetHitOfDailyByRange")
		return
	}

	for i, key := range keys {
		if v.([]interface{})[i] != nil {
			dailyValue, suberr := strconv.ParseInt(string(v.([]interface{})[i].([]byte)), 10, 64)
			if suberr != nil {
				err = errors.Wrap(suberr, "[err] GetHitOfDailyByRange")
			}
			scores = append(scores, &Score{Name: key.(string), Value: dailyValue})
		} else {
			scores = append(scores, nil)
		}
	}
	return
}
