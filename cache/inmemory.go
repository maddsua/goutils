package cache

import (
	"errors"
	"sync"
	"time"
)

type StorageTyper interface {
	StorageType() string
}

const NeverExpire = time.Duration(0)
const ExpireNow = time.Duration(-1)

type InMemoryCache struct {
	entries     map[string]*InmemoryEntry
	mtx         sync.Mutex
	nextCleanup time.Time
}

type InmemoryEntry struct {
	val any
	exp time.Time
}

func (this *InmemoryEntry) isExpired() bool {
	return !this.exp.IsZero() && this.exp.Before(time.Now())
}

func (this *InmemoryEntry) expire(ttl time.Duration) {

	if ttl <= ExpireNow {
		this.exp = time.Unix(1, 0)
		return
	}

	if ttl == NeverExpire {
		this.exp = time.Time{}
		return
	}

	this.exp = time.Now().Add(ttl)
}

func (this *InMemoryCache) subtleInit() {
	if this.entries == nil {
		this.entries = map[string]*InmemoryEntry{}
	}
}

func (this *InMemoryCache) StorageType() string {
	return "mem"
}

func (this *InMemoryCache) Set(key string, val any, ttl time.Duration) error {

	if ttl <= ExpireNow {
		return errors.New("cannot set an entry with ttl in past")
	}

	this.subtleInit()

	this.mtx.Lock()
	defer this.mtx.Unlock()

	now := time.Now()
	if now.After(this.nextCleanup) {
		go this.cleanupRoutine(now)
	}

	entry := &InmemoryEntry{
		val: val,
	}

	entry.expire(ttl)

	this.entries[key] = entry

	return nil
}

func (this *InMemoryCache) Get(key string) (any, bool) {

	this.subtleInit()

	this.mtx.Lock()
	defer this.mtx.Unlock()

	if entry, has := this.entries[key]; has && !entry.isExpired() {
		return entry.val, true
	}

	return nil, false
}

func (this *InMemoryCache) cleanupRoutine(now time.Time) {

	this.mtx.Lock()
	defer this.mtx.Unlock()

	for key, entry := range this.entries {
		if entry.isExpired() {
			delete(this.entries, key)
		}
	}

	this.nextCleanup = now.Add(time.Minute)
}

func (this *InMemoryCache) Expire(key string, ttl time.Duration) bool {

	this.subtleInit()

	this.mtx.Lock()
	defer this.mtx.Unlock()

	entry := this.entries[key]
	if entry == nil {
		return false
	}

	entry.expire(ttl)

	return true
}

func (this *InMemoryCache) TTL(key string) (time.Duration, bool) {

	this.subtleInit()

	this.mtx.Lock()
	defer this.mtx.Unlock()

	entry := this.entries[key]
	if entry == nil || entry.exp.IsZero() {
		return 0, false
	}

	return time.Until(entry.exp), true
}
