package cache_test

import (
	"testing"
	"time"

	"github.com/maddsua/golib/cache"
)

func TestInmemory_1(t *testing.T) {

	ic := cache.InMemoryCache{}

	ic.Set("goth", "nice", cache.NeverExpire)

	if val, _ := ic.Get("goth"); val.(string) != "nice" {
		t.Fatalf("invalid cached value: %v", val)
	}
}

func TestInmemory_2(t *testing.T) {

	ic := cache.InMemoryCache{}

	ic.Set("goth", "nice", time.Second)

	if val, _ := ic.Get("goth"); val.(string) != "nice" {
		t.Fatalf("invalid cached value: %v", val)
	}

	time.Sleep(time.Duration(1.2 * float64(time.Second)))

	if _, has := ic.Get("goth"); has {
		t.Fatalf("cached value was evicted")
	}
}
