package gocache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// NewRedisCache return a Cache backed by a Redis instance
func NewRedisCache(uri, password, prefix string) (Cache, error) {
	return &redisCache{
		redis: redis.NewClient(&redis.Options{
			Addr:     uri,
			Password: password,
			DB:       0,
		}),
		prefix: prefix,
	}, nil
}

type redisCache struct {
	redis  *redis.Client
	prefix string
}

func (rc *redisCache) Int64(key string) (int64, bool, error) {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	if val, err := rc.redis.Get(context.Background(), key).Int64(); err != nil {
		if errors.Is(err, redis.Nil) {
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

func (rc *redisCache) GetInt64(key string, callback func() (int64, time.Duration, error)) (int64, error) {
	val, exists, err := rc.Int64(key)
	if err != nil {
		return 0, err
	}

	if !exists {
		val, ttl, err := callback()
		if err != nil {
			return 0, err
		}

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
		if errors.Is(err, redis.Nil) {
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

func (rc *redisCache) GetInt(key string, callback func() (int, time.Duration, error)) (int, error) {
	val, exists, err := rc.Int(key)
	if err != nil {
		return 0, err
	}

	if !exists {
		val, ttl, err := callback()
		if err != nil {
			return 0, err
		}

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
		if errors.Is(err, redis.Nil) {
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

func (rc *redisCache) GetTime(key string, callback func() (time.Time, time.Duration, error)) (time.Time, error) {
	val, exists, err := rc.Time(key)
	if err != nil {
		return time.Time{}, err
	}

	if !exists {
		val, ttl, err := callback()
		if err != nil {
			return time.Time{}, err
		}

		if err := rc.SetTime(key, val, ttl); err != nil {
			return time.Time{}, err
		}

		return val, nil
	}

	return val, nil
}

func (rc *redisCache) Bytes(key string) ([]byte, bool, error) {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	if val, err := rc.redis.Get(context.Background(), key).Bytes(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}

		return nil, false, err
	} else {
		return val, true, nil
	}
}

func (rc *redisCache) SetBytes(key string, value []byte, ttl time.Duration) error {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	return rc.redis.Set(context.Background(), key, value, ttl).Err()
}

func (rc *redisCache) GetBytes(key string, callback func() ([]byte, time.Duration, error)) ([]byte, error) {
	val, exists, err := rc.Bytes(key)
	if err != nil {
		return nil, err
	}

	if !exists {
		val, ttl, err := callback()
		if err != nil {
			return nil, err
		}

		if err := rc.SetBytes(key, val, ttl); err != nil {
			return nil, err
		}

		return val, nil
	}

	return val, nil
}

func (rc *redisCache) Value(key string) (interface{}, bool, error) {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	b, err := rc.redis.Get(context.Background(), key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}

		return nil, false, err
	}

	var result interface{}
	msg := json.RawMessage(b)

	if err := json.Unmarshal(msg, &result); err != nil {
		return nil, false, err
	}

	return result, true, nil
}

func (rc *redisCache) SetValue(key string, value interface{}, ttl time.Duration) error {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rc.redis.Set(context.Background(), key, b, ttl).Err()
}

func (rc *redisCache) GetValue(key string, callback func() (interface{}, time.Duration, error)) (interface{}, error) {
	val, exists, err := rc.Value(key)
	if err != nil {
		return nil, err
	}

	if !exists {
		val, ttl, err := callback()
		if err != nil {
			return nil, err
		}

		if err := rc.SetValue(key, val, ttl); err != nil {
			return nil, err
		}

		return val, nil
	}

	return val, nil
}

func (rc *redisCache) String(key string) (string, bool, error) {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	s, err := rc.redis.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", false, nil
		}

		return "", false, err
	}

	return s, true, nil
}

func (rc *redisCache) SetString(key string, value string, ttl time.Duration) error {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	return rc.redis.Set(context.Background(), key, value, ttl).Err()
}

func (rc *redisCache) GetString(key string, callback func() (string, time.Duration, error)) (string, error) {
	val, exists, err := rc.String(key)
	if err != nil {
		return "", err
	}

	if !exists {
		val, ttl, err := callback()
		if err != nil {
			return "", err
		}

		if err := rc.SetString(key, val, ttl); err != nil {
			return "", err
		}

		return val, nil
	}

	return val, nil
}

func (rc *redisCache) Delete(key string) error {
	key = fmt.Sprintf("%s:%s", rc.prefix, key)

	err := rc.redis.Del(context.Background(), key).Err()
	if errors.Is(err, redis.Nil) {
		return nil
	}

	return err
}
