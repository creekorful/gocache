package gocache

import (
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

	// Bytes retrieve the []byte value identified by given key
	// it returns the value, a boolean indicating if the value exists or not, and an error (if any)
	Bytes(key string) ([]byte, bool, error)
	// SetBytes set the []byte value identified by given key
	// it takes the key, the value to be set, and an optional TTL (see: NoExpiration)
	SetBytes(key string, value []byte, ttl time.Duration) error
	// GetBytes either retrieve the value identified by given key if it exists
	// otherwise it will run the callback, cache and return the value
	GetBytes(key string, callback func() ([]byte, time.Duration)) ([]byte, error)

	// Value retrieve the raw value identified by given key
	// it returns the value, a boolean indicating if the value exists or not, and an error (if any)
	Value(key string) (interface{}, bool, error)
	// SetValue set the raw value identified by given key
	// it takes the key, the value to be set, and an optional TTL (see: NoExpiration)
	SetValue(key string, value interface{}, ttl time.Duration) error
	// GetValue either retrieve the raw value identified by given key if it exists
	// otherwise it will run the callback, cache and return the value
	GetValue(key string, callback func() (interface{}, time.Duration)) (interface{}, error)

	// Delete the value identified by given key.
	// the function does not fail if key does not exist
	Delete(key string) error
}
