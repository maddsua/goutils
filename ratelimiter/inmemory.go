package ratelimiter

import (
	"context"
	"time"

	"github.com/maddsua/goutils/cache"
)

func NewInmemory() *inmemory {
	return &inmemory{storage: &cache.InMemoryCache{}}
}

type inmemory struct {
	storage *cache.InMemoryCache
}

type inmemoryState struct {
	actions int
}

func (this *inmemory) Type() string {
	return "mem"
}

func (this *inmemory) mkKey(clientID string, actionID string) string {
	return clientID + ":" + actionID
}

func (this *inmemory) Use(clientID string, action Action) (Stats, error) {

	key := this.mkKey(clientID, action.ID)

	var state *inmemoryState
	if entry, _ := this.storage.Get(key); entry != nil {
		state = entry.(*inmemoryState)
		this.storage.Expire(key, action.Window)
	} else {
		state = &inmemoryState{}
		this.storage.Set(key, state, action.Window)
	}

	expires := time.Now().Add(action.Window)

	if state.actions >= MaxActionCount {
		return Stats{
			Quota:   0,
			Actions: state.actions,
			Expires: expires,
		}, nil
	}

	quoteActions := state.actions
	state.actions++

	return Stats{
		Quota:   clampQuota(action.Quota - quoteActions),
		Actions: state.actions,
		Expires: expires,
	}, nil
}

func (this *inmemory) UseContext(_ctx context.Context, clientID string, action Action) (Stats, error) {
	return this.Use(clientID, action)
}
