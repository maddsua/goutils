package ratelimiter

import (
	"context"
	"time"
)

type Ratelimter interface {
	Type() string
	Use(clientID string, action Action) (Stats, error)
	UseContext(ctx context.Context, clientID string, action Action) (Stats, error)
}

type Action struct {
	ID     string
	Quota  int
	Window time.Duration
}

type Stats struct {
	Quota   int
	Actions int
	Expires time.Time
}

func (this Stats) Exhausted() bool {
	return this.Quota <= 0
}

func clampQuota(quota int) int {
	if quota < 0 {
		return 0
	}
	return quota
}

const MaxActionCount = 1000000000
