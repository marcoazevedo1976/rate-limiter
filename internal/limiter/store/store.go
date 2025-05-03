package store

import "time"

type Store interface {
	Increment(key string, limit int, blockSeconds int) (allowed bool, ttl time.Duration, err error)
	IsBlocked(key string) (bool, time.Duration, error)
	GetRequestCount(key string) (int64, error)
}
