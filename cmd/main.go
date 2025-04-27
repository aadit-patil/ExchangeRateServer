package main

import (
	"log"
	"net/http"

	"github.com/aadit-patil/ExchangeRateServer/internal/api"
	"github.com/aadit-patil/ExchangeRateServer/internal/cache"
	"github.com/aadit-patil/ExchangeRateServer/internal/db"
	"github.com/aadit-patil/ExchangeRateServer/internal/metrics"
	model "github.com/aadit-patil/ExchangeRateServer/internal/models"
	"github.com/aadit-patil/ExchangeRateServer/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	//db.InitMySQL("root:@tcp(localhost:3306)/exchange")
	db.InitMySQL("root:root@tcp(db:3306)/exchange")
	cache.InitSingleton()

	strategy := service.NewCacheDBAPIStrategy()
	service.SetGlobalStrategy(strategy)

	// Prefetch DB data into cache
	model.PrefetchCache()

	http.HandleFunc("/convert", api.ConvertHandler)
	http.HandleFunc("/convert/range", api.ConvertRangeHandler)
	metrics.Init()
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Server running on :8088")

	log.Fatal(http.ListenAndServe(":8088", nil))
}
