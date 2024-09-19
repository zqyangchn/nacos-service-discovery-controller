package basicutils

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	*rate.Limiter
}

func InitRateLimiter(timeDuration string, b int) (*RateLimiter, error) {
	interval, err := time.ParseDuration(timeDuration)
	if err != nil {
		return nil, err
	}
	burstsEverySecond := b / int(interval.Seconds())

	return &RateLimiter{
		Limiter: rate.NewLimiter(rate.Limit(burstsEverySecond), b),
	}, nil
}

func (r *RateLimiter) WaitLimiter() {
	_ = r.Wait(context.Background())
}
