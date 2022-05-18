package gocache

import (
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

const (
	NoExpiration time.Duration = 0
)

// Cache represent a cache
type Cache interface {
	// Int64 retrieve the int64 value identified by given key
	// it returns the value, a boolean indicating if the value exists or not, and an error (if any)
	Int64(key string) (int64, bool, error)
	// SetInt64 set the int64 value identified by given key
	// it takes the key, the value to be set, and an optional TTL (see: NoExpiration)
	SetInt64(key string, value int64, ttl time.Duration) error
	// GetInt64 either retrieve the value identified by given key if it exists
	// otherwise it will run the callback, cache and return the value
	GetInt64(key string, callback func() (int64, time.Duration)) (int64, error)

	// Int retrieve the int value identified by given key
	// it returns the value, a boolean indicating if the value exists or not, and an error (if any)
	Int(key string) (int, bool, error)
	// SetInt set the int value identified by given key
	// it takes the key, the value to be set, and an optional TTL (see: NoExpiration)
	SetInt(key string, value int, ttl time.Duration) error
	// GetInt either retrieve the value identified by given key if it exists
	// otherwise it will run the callback, cache and return the value
	GetInt(key string, callback func() (int, time.Duration)) (int, error)

	// Time retrieve the time.Time value identified by given key
	// it returns the value, a boolean indicating if the value exists or not, and an error (if any)
	Time(key string) (time.Time, bool, error)
	// SetTime set the time.Time value identified by given key
	// it takes the key, the value to be set, and an optional TTL (see: NoExpiration)
	SetTime(key string, value time.Time, ttl time.Duration) error
	// GetTime either retrieve the value identified by given key if it exists
	// otherwise it will run the callback, cache and return the value
	GetTime(key string, callback func() (time.Time, time.Duration)) (time.Time, error)

	// Delete the value identified by given key.
	// the function does not fail if key does not exist
	Delete(key string) error
}

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

// NewMemoryCache return a Cache in-memory. (Slow)
func NewMemoryCache(prefix string) Cache {
	return &memoryCache{
		values: map[string]entry{},
		mutex:  sync.RWMutex{},
		prefix: prefix,
	}
}
