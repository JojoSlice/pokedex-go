package internal

import (
	"testing"
	"time"
)

func TestReaploop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 7*time.Millisecond
	entry := cacheEntry{time.Now(), []byte("testdata")}

	cache := NewCache(baseTime)
	cache.Set("https://testing.dev", entry)

	_, ok := cache.Get("https://testing.dev")
	if !ok {
		t.Fatal("expected entry to exist before reap")
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://testing.dev")
	if ok {
		t.Fatal("expected entry to be reaped, but it still exists")
	}
}
