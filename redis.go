package gocache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisCache struct {
	redis  *redis.Client
	prefix string
}

func (rc *redisCache) Int64(key string) (int64, bool, error) {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	if val, err := rc.redis.Get(context.Background(), key).Int64(); err != nil {
		if err == redis.Nil {
			return 0, false, nil
		}

		return 0, false, err
	} else {
		return val, true, nil
	}
}

func (rc *redisCache) SetInt64(key string, value int64, ttl time.Duration) error {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	return rc.redis.Set(context.Background(), key, value, ttl).Err()
}

func (rc *redisCache) GetInt64(key string, callback func() (int64, time.Duration)) (int64, error) {
	val, exists, err := rc.Int64(key)
	if err != nil {
		return 0, err
	}

	if !exists {
		val, ttl := callback()
		if err := rc.SetInt64(key, val, ttl); err != nil {
			return 0, err
		}

		return val, nil
	}

	return val, nil
}

func (rc *redisCache) Int(key string) (int, bool, error) {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	if val, err := rc.redis.Get(context.Background(), key).Int(); err != nil {
		if err == redis.Nil {
			return 0, false, nil
		}

		return 0, false, err
	} else {
		return val, true, nil
	}
}

func (rc *redisCache) SetInt(key string, value int, ttl time.Duration) error {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	return rc.redis.Set(context.Background(), key, value, ttl).Err()
}

func (rc *redisCache) GetInt(key string, callback func() (int, time.Duration)) (int, error) {
	val, exists, err := rc.Int(key)
	if err != nil {
		return 0, err
	}

	if !exists {
		val, ttl := callback()
		if err := rc.SetInt(key, val, ttl); err != nil {
			return 0, err
		}

		return val, nil
	}

	return val, nil
}

func (rc *redisCache) Time(key string) (time.Time, bool, error) {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	if val, err := rc.redis.Get(context.Background(), key).Time(); err != nil {
		if err == redis.Nil {
			return time.Time{}, false, nil
		}

		return time.Time{}, false, err
	} else {
		return val, true, nil
	}
}

func (rc *redisCache) SetTime(key string, value time.Time, ttl time.Duration) error {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	return rc.redis.Set(context.Background(), key, value, ttl).Err()
}

func (rc *redisCache) GetTime(key string, callback func() (time.Time, time.Duration)) (time.Time, error) {
	val, exists, err := rc.Time(key)
	if err != nil {
		return time.Time{}, err
	}

	if !exists {
		val, ttl := callback()
		if err := rc.SetTime(key, val, ttl); err != nil {
			return time.Time{}, err
		}

		return val, nil
	}

	return val, nil
}

func (rc *redisCache) Delete(key string) error {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	err := rc.redis.Del(context.Background(), key).Err()
	if err == redis.Nil {
		return nil
	}

	return err
}
