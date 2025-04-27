package service

import (
	"testing"
	"time"

	"github.com/aadit-patil/ExchangeRateServer/internal/cache"
	"github.com/aadit-patil/ExchangeRateServer/internal/db"
	"github.com/aadit-patil/ExchangeRateServer/internal/mocks"
	"github.com/golang/mock/gomock"
)

func TestGetRate_FromCache(t *testing.T) {
	cache.InitSingleton()
	key := time.Now().Format("2006-01-02") + ":USD"
	cache.SetRate(key, 83.0, time.Hour)

	strategy := NewCacheDBAPIStrategy()
	SetGlobalStrategy(strategy)

	rate, err := strategy.GetRate("USD", "INR", "")
	if err != nil || rate != 83.0 {
		t.Errorf("expected 83.0 from cache, got %v, err: %v", rate, err)
	}
}

func TestGetRate_FromDB(t *testing.T) {
	cache.InitSingleton()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	db.SetDatabase(mockDB)

	date := time.Now().Format("2006-01-02")
	mockDB.EXPECT().GetRate("USD", "EUR", date).Return(0.91, nil)

	strategy := NewCacheDBAPIStrategy()
	SetGlobalStrategy(strategy)

	rate, err := strategy.GetRate("USD", "EUR", date)
	if err != nil || rate != 0.91 {
		t.Errorf("expected 0.91 from DB, got %v, err: %v", rate, err)
	}
}
