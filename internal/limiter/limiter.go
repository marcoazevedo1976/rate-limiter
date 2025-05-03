package limiter

import (
	"net/http"

	"rate-limiter/internal/config"
	"rate-limiter/internal/limiter/store"
)

type RateLimiter struct {
	cfg   *config.Config
	store store.Store
}

func NewRateLimiter(cfg *config.Config, store store.Store) *RateLimiter {
	return &RateLimiter{
		cfg:   cfg,
		store: store,
	}
}

func (r *RateLimiter) Allow(key string, isToken bool) (bool, int, string) {
	var limit, block int
	if isToken {
		limit = r.cfg.RateLimitToken
		block = r.cfg.BlockDurationToken
	} else {
		limit = r.cfg.RateLimitIP
		block = r.cfg.BlockDurationIP
	}

	blocked, _, err := r.store.IsBlocked(key)
	if err != nil {
		return false, http.StatusInternalServerError, "internal error"
	}
	if blocked {
		return false, http.StatusTooManyRequests, "you have reached the maximum number of requests or actions allowed within a certain time frame"
	}

	allowed, _, err := r.store.Increment(key, limit, block)
	if err != nil {
		return false, http.StatusInternalServerError, "internal error"
	}
	if !allowed {
		return false, http.StatusTooManyRequests, "you have reached the maximum number of requests or actions allowed within a certain time frame"
	}

	return true, http.StatusOK, ""
}
