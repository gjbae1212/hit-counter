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

func (d *db) IncreaseHitOfDaily(id string) (*Score, error) {
	if id == "" {
		return nil, fmt.Errorf("[err] IncreaseHitOfDaily empty param")
	}

	daily := allan_util.TimeToDailyStringFormat(time.Now())
	key := fmt.Sprintf(hitDailyFormat, daily, id)
	v, err := d.redis.DoWithTimeout(timeout, "INCR", key)
	if err != nil {
		return nil, errors.Wrap(err, "[err] IncreaseHitOfDaily")
	}

	return &Score{Name: id, Value: v.(int64)}, nil
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
	return &Score{Name: id, Value: v.(int64)}, nil
}

func (d *db) GetHitOfDaily(id string) (*Score, error) {
	if id == "" {
		return nil, fmt.Errorf("[err] GetHitOfDaily empty param")
	}

	daily := allan_util.TimeToDailyStringFormat(time.Now())
	key := fmt.Sprintf(hitDailyFormat, daily, id)

	v, err := d.redis.DoWithTimeout(timeout, "GET", key)
	if err != nil {
		return nil, errors.Wrap(err, "[err] GetHitOfDaily")
	}

	rt, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "[err] GetHitOfDaily")
	}

	return &Score{Name: id, Value: rt}, nil
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

	rt, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "[err] GetHitOfTotal")
	}

	return &Score{Name: id, Value: rt}, nil
}
