package cache

import (
	"log"
	"time"

	"github.com/dgraph-io/ristretto"
)

var instance *ristretto.Cache

func InitSingleton() {
	var err error
	instance, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,
		MaxCost:     1 << 28,
		BufferItems: 128,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func GetCache() *ristretto.Cache {
	return instance
}

func SetRate(key string, rate float64, ttl time.Duration) {
	instance.SetWithTTL(key, rate, 1, ttl)
}

func GetRate(key string) (float64, bool) {
	val, found := instance.Get(key)
	if !found {
		return 0, false
	}
	return val.(float64), true
}
