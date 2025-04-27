package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/aadit-patil/ExchangeRateServer/internal/cache"
	"github.com/aadit-patil/ExchangeRateServer/internal/client"
	"github.com/aadit-patil/ExchangeRateServer/internal/db"
)

type RateFetchStrategy interface {
	GetRate(from, to, date string) (float64, error)
	FetchAndStoreRate(from, to, date string)
}

var strategy RateFetchStrategy

func SetGlobalStrategy(s RateFetchStrategy) {
	strategy = s
}

func GetGlobalStrategy() RateFetchStrategy {
	return strategy
}

type CacheDBAPIStrategy struct{}

func NewCacheDBAPIStrategy() *CacheDBAPIStrategy {
	return &CacheDBAPIStrategy{}
}

func (s *CacheDBAPIStrategy) GetRate(from, to, date string) (float64, error) {
	key := fmt.Sprintf("%s:%s:%s", date, from, to)
	if rate, ok := cache.GetRate(key); ok {
		return rate, nil
	}

	rate, err := db.DBImpl.GetRate(from, to, date)
	if err == nil {
		cache.SetRate(key, rate, time.Hour)
		return rate, nil
	}

	return 0, errors.New("rate not available")
}

func (s *CacheDBAPIStrategy) FetchAndStoreRate(from, to, date string) {
	rates, err := client.FetchRatesForBase(from, date)
	if err == nil {

		_ = db.DBImpl.InsertMultipleRates(from, date, rates, time.Now().Add(1*time.Hour))

		key := fmt.Sprintf("%s:%s:%s", date, from, to)
		cache.SetRate(key, rates[to], time.Hour)

	}
}
