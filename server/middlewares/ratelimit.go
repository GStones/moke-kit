package middlewares

import (
	"context"
	"errors"

	"golang.org/x/time/rate"
)

var (
	ErrRateLimit = errors.New("rate limit exceeded")
)

// RateLimiter rate limit
type RateLimiter struct {
	tokenLimiter *rate.Limiter
}

// CreateRateLimiter creates a rate limiter
// here we use golang.org/x/time/rate
func CreateRateLimiter(num int) *RateLimiter {
	tokenLimiter := rate.NewLimiter(rate.Limit(num), num)
	return &RateLimiter{
		tokenLimiter: tokenLimiter,
	}
}

func (rl *RateLimiter) Limit(_ context.Context) error {
	if !rl.tokenLimiter.Allow() {
		return ErrRateLimit
	}
	return nil
}
