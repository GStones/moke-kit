package rate

import (
	"context"

	"go.uber.org/ratelimit"
)

type RateLimiter struct {
	limiter ratelimit.Limiter
}

func CreateRateLimiter(rate int) *RateLimiter {
	rl := ratelimit.New(rate) // per second, some slack.
	return &RateLimiter{
		limiter: rl,
	}
}

func (rl *RateLimiter) Limit(_ context.Context) error {
	rl.limiter.Take()
	return nil
}
