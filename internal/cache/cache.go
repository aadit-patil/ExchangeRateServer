package cache

import (
	"sync"
	"time"
)

type entry struct {
	Rate      float64
	Timestamp time.Time
}

var (
	once     sync.Once
	instance *Cache
)

type Cache struct {
	mu    sync.RWMutex
	store map[string]entry
}

func InitSingleton() {
	once.Do(func() {
		instance = &Cache{store: make(map[string]entry)}
	})
}

func GetCache() *Cache {
	return instance
}

func (c *Cache) Get(key string) (float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.store[key]
	if !ok || time.Since(val.Timestamp) > time.Hour {
		return 0, false
	}
	return val.Rate, true
}

func (c *Cache) Set(key string, rate float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = entry{Rate: rate, Timestamp: time.Now()}
}
