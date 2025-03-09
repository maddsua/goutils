package cache

import (
	"sync"
	"time"
)

const NeverExpire = time.Duration(0)

type InMemoryCache struct {
	entries     map[string]*InmemoryEntry
	mtx         sync.Mutex
	nextCleanup time.Time
}

type InmemoryEntry struct {
	Val any
	Exp time.Time
}

func (this *InmemoryEntry) IsExpired() bool {
	return !this.Exp.IsZero() && this.Exp.Before(time.Now())
}

func (this *InMemoryCache) subtleInit() {
	if this.entries == nil {
		this.entries = map[string]*InmemoryEntry{}
	}
}

func (this *InMemoryCache) Type() string {
	return "mem"
}

func (this *InMemoryCache) Set(key string, val any, ttl time.Duration) {

	this.subtleInit()

	this.mtx.Lock()
	defer this.mtx.Unlock()

	now := time.Now()
	if now.After(this.nextCleanup) {
		go this.cleanupRoutine(now)
	}

	var expires time.Time
	if ttl > NeverExpire {
		expires = time.Now().Add(ttl)
	}

	this.entries[key] = &InmemoryEntry{
		Val: val,
		Exp: expires,
	}
}

func (this *InMemoryCache) Get(key string) (any, bool) {

	this.subtleInit()

	this.mtx.Lock()
	defer this.mtx.Unlock()

	if entry, has := this.entries[key]; has && !entry.IsExpired() {
		return entry.Val, true
	}

	return nil, false
}

func (this *InMemoryCache) GetEntry(key string) *InmemoryEntry {

	this.subtleInit()

	this.mtx.Lock()
	defer this.mtx.Unlock()

	return this.entries[key]
}

func (this *InMemoryCache) cleanupRoutine(now time.Time) {

	this.mtx.Lock()
	defer this.mtx.Unlock()

	for key, entry := range this.entries {
		if entry.IsExpired() {
			delete(this.entries, key)
		}
	}

	this.nextCleanup = now.Add(time.Minute)
}
