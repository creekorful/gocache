package gocache

import (
	"fmt"
	"sync"
	"time"
)

type entry struct {
	value          interface{}
	expirationTime time.Time
}

// NewMemoryCache return a Cache in-memory.
func NewMemoryCache(prefix string) Cache {
	return &memoryCache{
		values: map[string]entry{},
		mutex:  sync.RWMutex{},
		prefix: prefix,
	}
}

type memoryCache struct {
	values map[string]entry
	mutex  sync.RWMutex
	prefix string
}

func (mc *memoryCache) Int64(key string) (int64, bool, error) {
	val, exists, err := mc.Value(key)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}

	return val.(int64), true, nil
}

func (mc *memoryCache) SetInt64(key string, value int64, ttl time.Duration) error {
	return mc.SetValue(key, value, ttl)
}

func (mc *memoryCache) GetInt64(key string, callback func() (int64, time.Duration, error)) (int64, error) {
	val, exists, err := mc.Int64(key)
	if err != nil {
		return 0, err
	}
	if !exists {
		val, ttl, err := callback()
		if err != nil {
			return 0, err
		}

		if err := mc.SetValue(key, val, ttl); err != nil {
			return 0, err
		}

		return val, nil
	}

	return val, nil
}

func (mc *memoryCache) Int(key string) (int, bool, error) {
	val, exists, err := mc.Value(key)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}

	return val.(int), true, nil
}

func (mc *memoryCache) SetInt(key string, value int, ttl time.Duration) error {
	return mc.SetValue(key, value, ttl)
}

func (mc *memoryCache) GetInt(key string, callback func() (int, time.Duration, error)) (int, error) {
	val, exists, err := mc.Int(key)
	if err != nil {
		return 0, err
	}
	if !exists {
		val, ttl, err := callback()
		if err != nil {
			return 0, err
		}

		if err := mc.SetValue(key, val, ttl); err != nil {
			return 0, err
		}

		return val, nil
	}

	return val, nil
}

func (mc *memoryCache) Time(key string) (time.Time, bool, error) {
	val, exists, err := mc.Value(key)
	if err != nil {
		return time.Time{}, false, err
	}
	if !exists {
		return time.Time{}, false, nil
	}

	return val.(time.Time), true, nil
}

func (mc *memoryCache) SetTime(key string, value time.Time, ttl time.Duration) error {
	return mc.SetValue(key, value, ttl)
}

func (mc *memoryCache) GetTime(key string, callback func() (time.Time, time.Duration, error)) (time.Time, error) {
	val, exists, err := mc.Time(key)
	if err != nil {
		return time.Time{}, err
	}
	if !exists {
		val, ttl, err := callback()
		if err != nil {
			return time.Time{}, err
		}

		if err := mc.SetValue(key, val, ttl); err != nil {
			return time.Time{}, err
		}

		return val, nil
	}

	return val, nil
}

func (mc *memoryCache) Bytes(key string) ([]byte, bool, error) {
	val, exists, err := mc.Value(key)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}

	return val.([]byte), true, nil
}

func (mc *memoryCache) SetBytes(key string, value []byte, ttl time.Duration) error {
	return mc.SetValue(key, value, ttl)
}

func (mc *memoryCache) GetBytes(key string, callback func() ([]byte, time.Duration, error)) ([]byte, error) {
	val, exists, err := mc.Bytes(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		val, ttl, err := callback()
		if err != nil {
			return nil, err
		}

		if err := mc.SetValue(key, val, ttl); err != nil {
			return nil, err
		}

		return val, nil
	}

	return val, nil
}

func (mc *memoryCache) Value(key string) (interface{}, bool, error) {
	if mc.prefix != "" {
		key = fmt.Sprintf("%s:%s", mc.prefix, key)
	}

	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	val, exists := mc.values[key]
	if !exists {
		return nil, false, nil
	}

	// handle expiration
	if !val.expirationTime.IsZero() && val.expirationTime.Before(time.Now()) {
		delete(mc.values, key)
		return nil, false, nil
	}

	return val.value, true, nil
}

func (mc *memoryCache) SetValue(key string, value interface{}, ttl time.Duration) error {
	if mc.prefix != "" {
		key = fmt.Sprintf("%s:%s", mc.prefix, key)
	}

	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	expirationTime := time.Time{}
	if ttl != NoExpiration {
		expirationTime = time.Now().Add(ttl)
	}

	mc.values[key] = entry{
		value:          value,
		expirationTime: expirationTime,
	}

	return nil
}

func (mc *memoryCache) GetValue(key string, callback func() (interface{}, time.Duration, error)) (interface{}, error) {
	val, exists, err := mc.Value(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		val, ttl, err := callback()
		if err != nil {
			return nil, err
		}

		if err := mc.SetValue(key, val, ttl); err != nil {
			return nil, err
		}

		return val, nil
	}

	return val, nil
}

func (mc *memoryCache) String(key string) (string, bool, error) {
	if mc.prefix != "" {
		key = fmt.Sprintf("%s:%s", mc.prefix, key)
	}

	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	val, exists := mc.values[key]
	if !exists {
		return "", false, nil
	}

	// handle expiration
	if !val.expirationTime.IsZero() && val.expirationTime.Before(time.Now()) {
		delete(mc.values, key)
		return "", false, nil
	}

	return val.value.(string), true, nil
}

func (mc *memoryCache) SetString(key string, value string, ttl time.Duration) error {
	if mc.prefix != "" {
		key = fmt.Sprintf("%s:%s", mc.prefix, key)
	}

	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	expirationTime := time.Time{}
	if ttl != NoExpiration {
		expirationTime = time.Now().Add(ttl)
	}

	mc.values[key] = entry{
		value:          value,
		expirationTime: expirationTime,
	}

	return nil
}

func (mc *memoryCache) GetString(key string, callback func() (string, time.Duration, error)) (string, error) {
	val, exists, err := mc.String(key)
	if err != nil {
		return "", err
	}
	if !exists {
		val, ttl, err := callback()
		if err != nil {
			return "", err
		}

		if err := mc.SetString(key, val, ttl); err != nil {
			return "", err
		}

		return val, nil
	}

	return val, nil
}

func (mc *memoryCache) Delete(key string) error {
	if mc.prefix != "" {
		key = fmt.Sprintf("%s:%s", mc.prefix, key)
	}

	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	delete(mc.values, key)

	return nil
}
