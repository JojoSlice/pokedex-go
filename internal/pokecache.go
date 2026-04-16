package internal

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	data     map[string]cacheEntry
	mu       sync.RWMutex
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		interval: interval,
		data:     make(map[string]cacheEntry),
	}
	go c.reapLoop()
	return c
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.data {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) Set(key string, value cacheEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

func (c *Cache) Get(key string) (cacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	return value, ok
}
