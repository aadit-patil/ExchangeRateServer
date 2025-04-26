package main

import (
	"log"
	"net/http"

	"github.com/aadit-patil/ExchangeRateServer/internal/api"
	"github.com/aadit-patil/ExchangeRateServer/internal/cache"
	"github.com/aadit-patil/ExchangeRateServer/internal/db"
	"github.com/aadit-patil/ExchangeRateServer/internal/metrics"
	"github.com/aadit-patil/ExchangeRateServer/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	db.InitMySQL("root:@tcp(localhost:3306)/exchange")
	cache.InitSingleton()

	strategy := service.NewCacheDBAPIStrategy()
	service.SetGlobalStrategy(strategy)

	http.HandleFunc("/convert", api.ConvertHandler)
	log.Println("Server running on :8088")
	metrics.Init()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8088", nil))
}
