package main

import (
	"sync"
	"time"
)

type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
	stop  chan struct{}
}

type CacheItem struct {
	value    interface{}
	expireAt time.Time
}

func NewCache(cleanupInterval time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]CacheItem),
		stop:  make(chan struct{}),
	}

	// Start background cleaner
	if cleanupInterval > 0 {
		go c.cleanup(cleanupInterval)
	}
	return c
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().After(item.expireAt) {
		return nil, false
	}

	return item.value, true
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		value:    value,
		expireAt: time.Now().Add(ttl),
	}
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, item := range c.items {
				if time.Now().After(item.expireAt) {
					delete(c.items, key)
				}
			}
			c.mu.Unlock()
		case <-c.stop:
			return
		}
	}
}

func (c *Cache) Stop() {
	close(c.stop)
}
