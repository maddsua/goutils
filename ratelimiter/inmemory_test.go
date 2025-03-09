package ratelimiter_test

import (
	"testing"
	"time"

	"github.com/maddsua/goutils/ratelimiter"
)

func TestInmem_1(t *testing.T) {

	action := ratelimiter.Action{ID: "signup", Quota: 2, Window: time.Second}

	rl := ratelimiter.NewInmemory()

	if stats, err := rl.Use("localhost", action); err != nil {
		t.Fatalf("unexpected rl error (1): %v", err)
	} else if stats.Quota != 2 {
		t.Fatalf("expected quota to update (1): %v", stats)
	}

	if stats, err := rl.Use("localhost", action); err != nil {
		t.Fatalf("unexpected rl error (2): %v", err)
	} else if stats.Quota != 1 {
		t.Fatalf("expected quota to update (2): %v", stats)
	}

	if stats, err := rl.Use("localhost", action); err != nil {
		t.Fatalf("unexpected rl error (3): %v", err)
	} else if !stats.Exhausted() {
		t.Fatalf("expected quota to update (3): %v", stats)
	}

	if stats, err := rl.Use("localhost-2", action); err != nil {
		t.Fatalf("unexpected rl error (4): %v", err)
	} else if stats.Quota != 2 {
		t.Fatalf("unecpected quota rejection (4): %v", stats)
	}

	time.Sleep(time.Duration(1.2 * float64(time.Second)))

	if stats, err := rl.Use("localhost", action); err != nil {
		t.Fatalf("unexpected rl error (5): %v", err)
	} else if stats.Exhausted() {
		t.Fatalf("expected quota to reset (5): %v", stats)
	}
}
