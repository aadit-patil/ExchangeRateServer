package main

import (
	"log"
	"net/http"

	"github.com/aadit-patil/ExchangeRateServer/internal/api"   // Import the db package
	"github.com/aadit-patil/ExchangeRateServer/internal/cache" // Import the cache package
	"github.com/aadit-patil/ExchangeRateServer/internal/db"
	"github.com/aadit-patil/ExchangeRateServer/internal/service"
)

func main() {
	db.InitMySQL("root:root@tcp(db:3306)/exchange")
	cache.InitSingleton()

	strategy := service.NewCacheDBAPIStrategy()
	service.SetGlobalStrategy(strategy)

	http.HandleFunc("/convert", api.ConvertHandler)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
