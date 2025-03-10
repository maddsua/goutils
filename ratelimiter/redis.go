package ratelimiter

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(client RedisThinClient) *redislimiter {
	return &redislimiter{redis: client}
}

type RedisThinClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
}

type redislimiter struct {
	redis RedisThinClient
}

func (this *redislimiter) Type() string {
	return "redis"
}

func (this *redislimiter) UseContext(ctx context.Context, clientID string, action Action) (Stats, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	entryKey := "rlq:" + clientID + ":" + action.ID

	var counter int

	//	not using incr to avoid having to use expire command
	if entry, err := this.redis.Get(ctx, entryKey).Result(); err != nil {
		if err != redis.Nil {
			return Stats{}, err
		}
	} else if counter, err = strconv.Atoi(entry); err != nil {
		return Stats{}, err
	}

	if counter >= MaxActionCount {

		if err := this.redis.Expire(ctx, entryKey, action.Window).Err(); err != nil {
			return Stats{}, err
		}

		return Stats{
			Quota:   0,
			Actions: counter,
			Expires: time.Now().Add(action.Window),
		}, nil
	}

	counter++

	if err := this.redis.Set(ctx, entryKey, counter, action.Window).Err(); err != nil {
		return Stats{}, err
	}

	return Stats{
		Quota:   clampQuota(action.Quota - (counter - 1)),
		Actions: counter,
		Expires: time.Now().Add(action.Window),
	}, nil
}

func (this *redislimiter) Use(clientID string, action Action) (Stats, error) {
	return this.UseContext(context.Background(), clientID, action)
}
