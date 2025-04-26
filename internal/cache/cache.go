package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Rates     map[string]float64
	ExpiresAt time.Time
}

var (
	once     sync.Once
	instance *Cache
)

type Cache struct {
	mu    sync.RWMutex
	store map[string]CacheEntry
}

func InitSingleton() {
	once.Do(func() {
		instance = &Cache{
			store: make(map[string]CacheEntry),
		}
	})
}

func GetCache() *Cache {
	return instance
}

func (c *Cache) GetRates(baseDate string) (map[string]float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, exists := c.store[baseDate]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Rates, true
}

func (c *Cache) SetRates(baseDate string, newRates map[string]float64, ttl time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	existing, ok := c.store[baseDate]
	if ok && time.Now().Before(existing.ExpiresAt) {
		for k, v := range newRates {
			existing.Rates[k] = v
		}
		existing.ExpiresAt = ttl
		c.store[baseDate] = existing
	} else {
		c.store[baseDate] = CacheEntry{
			Rates:     newRates,
			ExpiresAt: ttl,
		}
	}
}
