package main

import (
	"log"
	"net/http"

	"github.com/aadit-patil/ExchangeRateServer/internal/api"
	"github.com/aadit-patil/ExchangeRateServer/internal/cache"
	"github.com/aadit-patil/ExchangeRateServer/internal/db"
	"github.com/aadit-patil/ExchangeRateServer/internal/service"
)

func main() {
	db.InitMySQL("root:@tcp(localhost:3306)/exchange")
	cache.InitSingleton()

	strategy := service.NewCacheDBAPIStrategy()
	service.SetGlobalStrategy(strategy)

	http.HandleFunc("/convert", api.ConvertHandler)
	log.Println("Server running on :8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
