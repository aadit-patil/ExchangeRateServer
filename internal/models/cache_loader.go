package model

import (
	"fmt"
	"time"

	"github.com/aadit-patil/ExchangeRateServer/internal/cache"
	"github.com/aadit-patil/ExchangeRateServer/internal/configs"
	"github.com/aadit-patil/ExchangeRateServer/internal/db"
	"github.com/aadit-patil/ExchangeRateServer/internal/service"
)

func PrefetchCache() {
	today := time.Now()
	startDate := today.AddDate(0, 0, -90)
	prefetched := false
	for d := startDate; !d.After(today); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")

		for from := range configs.SupportedCurrencies {
			for to := range configs.SupportedCurrencies {
				if from == to {
					continue
				}
				rate, err := db.DBImpl.GetRate(from, to, dateStr)
				if err == nil {
					key := fmt.Sprintf("%s:%s:%s", dateStr, from, to)
					cache.SetRate(key, rate, time.Hour)
				}
				if !prefetched {
					prefetched = true
					key := fmt.Sprintf("%s:%s:%s", today.Format("2006-01-02"), from, to)

					if _, ok := cache.GetRate(key); !ok {
						service.GetGlobalStrategy().FetchAndStoreRate(from, to, today.Format("2006-01-02"))
					}
				}
			}
		}
	}

}
