package service

import (
	"errors"

	"fmt"

	"github.com/aadit-patil/ExchangeRateServer/internal/cache"
	"github.com/aadit-patil/ExchangeRateServer/internal/client"
	"github.com/aadit-patil/ExchangeRateServer/internal/db"
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
	key := fmt.Sprintf("%s:%s:%s", date, from, to)
	if rate, ok := cache.GetCache().Get(key); ok {
		return rate, nil
	}

	rate, err := db.DBImpl.GetRate(from, to, date)
	if err == nil {
		cache.GetCache().Set(key, rate)
		return rate, nil
	}

	rate, err = client.FetchRate(from, to, date)
	if err != nil {
		return 0, errors.New("failed to retrieve rate from all sources")
	}
	db.DBImpl.InsertRate(from, to, date, rate)
	cache.GetCache().Set(key, rate)
	return rate, nil
}
