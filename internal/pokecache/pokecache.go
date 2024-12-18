package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]cacheEntry
	interval time.Duration
	mu       sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
		mu:       sync.Mutex{},
	}

	go cache.reapLoop()

	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return value.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key := range c.entries {
				timeSinceCreation := time.Since(c.entries[key].createdAt)
				if timeSinceCreation > c.interval {
					delete(c.entries, key)
				}
			}
			c.mu.Unlock()
		}
	}
}

// Use for troubleshooting
func (c *Cache) PrintCache() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.entries {
		fmt.Printf("Key: %s, Value: %s\n", k, string(v.val))
	}
}
