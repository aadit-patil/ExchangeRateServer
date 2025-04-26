package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/aadit-patil/ExchangeRateServer/internal/cache"
	"github.com/aadit-patil/ExchangeRateServer/internal/client"
	"github.com/aadit-patil/ExchangeRateServer/internal/db"
	"github.com/aadit-patil/ExchangeRateServer/internal/metrics"
)

type RateFetchStrategy interface {
	GetRate(from, to, date string) (float64, error)
}

var strategy RateFetchStrategy

func SetGlobalStrategy(s RateFetchStrategy) {
	strategy = s
}

type CacheDBAPIStrategy struct{}

func NewCacheDBAPIStrategy() *CacheDBAPIStrategy {
	return &CacheDBAPIStrategy{}
}

func (s *CacheDBAPIStrategy) GetRate(from, to, date string) (float64, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	key := fmt.Sprintf("%s:%s", date, from)

	if rates, ok := cache.GetCache().GetRates(key); ok {
		if rate, exists := rates[to]; exists {
			metrics.CacheHits.Inc()
			return rate, nil
		}
	}
	metrics.CacheMisses.Inc()

	rate, err := db.DBImpl.GetRate(from, to, date)
	if err == nil {
		cache.GetCache().SetRates(key, map[string]float64{to: rate}, time.Now().Add(1*time.Hour))
		return rate, nil
	}

	ratesMap, ttl, err := client.FetchRatesForBase(from)
	if err != nil {
		return 0, errors.New("failed to fetch from API")
	}

	errDB := db.DBImpl.InsertMultipleRates(from, date, ratesMap, ttl)
	if errDB != nil {
		return 0, fmt.Errorf("db get error: %w", err)
	}
	if rate, ok := ratesMap[to]; ok && supportedCurrencies[to] {
		cache.GetCache().SetRates(key, map[string]float64{to: rate}, ttl)
		return rate, nil
	}
	return 0, errors.New("unsupported currency or missing rate")
}
