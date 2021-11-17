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

type memoryCache struct {
	values map[string]entry
	mutex  sync.RWMutex
	prefix string
}

func (mc *memoryCache) Int64(key string) (int64, bool, error) {
	val, exists, err := mc.getValue(key)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}

	return val.(int64), true, nil
}

func (mc *memoryCache) SetInt64(key string, value int64, ttl time.Duration) error {
	return mc.setValue(key, value, ttl)
}

func (mc *memoryCache) GetInt64(key string, callback func() (int64, time.Duration)) (int64, error) {
	val, exists, err := mc.Int64(key)
	if err != nil {
		return 0, err
	}
	if !exists {
		val, ttl := callback()
		if err := mc.setValue(key, val, ttl); err != nil {
			return 0, err
		}

		return val, nil
	}

	return val, nil
}

func (mc *memoryCache) Int(key string) (int, bool, error) {
	val, exists, err := mc.getValue(key)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}

	return val.(int), true, nil
}

func (mc *memoryCache) SetInt(key string, value int, ttl time.Duration) error {
	return mc.setValue(key, value, ttl)
}

func (mc *memoryCache) GetInt(key string, callback func() (int, time.Duration)) (int, error) {
	val, exists, err := mc.Int(key)
	if err != nil {
		return 0, err
	}
	if !exists {
		val, ttl := callback()
		if err := mc.setValue(key, val, ttl); err != nil {
			return 0, err
		}

		return val, nil
	}

	return val, nil
}

func (mc *memoryCache) Time(key string) (time.Time, bool, error) {
	val, exists, err := mc.getValue(key)
	if err != nil {
		return time.Time{}, false, err
	}
	if !exists {
		return time.Time{}, false, nil
	}

	return val.(time.Time), true, nil
}

func (mc *memoryCache) SetTime(key string, value time.Time, ttl time.Duration) error {
	return mc.setValue(key, value, ttl)
}

func (mc *memoryCache) GetTime(key string, callback func() (time.Time, time.Duration)) (time.Time, error) {
	val, exists, err := mc.Time(key)
	if err != nil {
		return time.Time{}, err
	}
	if !exists {
		val, ttl := callback()
		if err := mc.setValue(key, val, ttl); err != nil {
			return time.Time{}, err
		}

		return val, nil
	}

	return val, nil
}

func (mc *memoryCache) Delete(key string) error {
	key = fmt.Sprintf("%s:%s", mc.prefix, key)

	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	delete(mc.values, key)

	return nil
}

func (mc *memoryCache) getValue(key string) (interface{}, bool, error) {
	key = fmt.Sprintf("%s:%s", mc.prefix, key)

	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	val, exists := mc.values[key]
	if !exists {
		return 0, false, nil
	}

	// handle expiration
	if !val.expirationTime.IsZero() && val.expirationTime.Before(time.Now()) {
		delete(mc.values, key)
		return 0, false, nil
	}

	return val.value, true, nil
}

func (mc *memoryCache) setValue(key string, value interface{}, ttl time.Duration) error {
	key = fmt.Sprintf("%s:%s", mc.prefix, key)

	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.values[key] = entry{
		value:          value,
		expirationTime: time.Now().Add(ttl),
	}

	return nil
}
