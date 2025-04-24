package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var supported = map[string]bool{"USD": true, "INR": true, "EUR": true, "JPY": true, "GBP": true}

type RateCache struct {
	mu        sync.RWMutex
	base      string
	rates     map[string]float64
	lastFetch time.Time
	apiKey    string
}

func NewRateCache(apiKey string, base string) *RateCache {
	return &RateCache{
		base:   base,
		apiKey: apiKey,
		rates:  make(map[string]float64),
	}
}

// Response struct for exchangerate-api.com
type ExchangeRateAPIResponse struct {
	Result          string             `json:"result"`
	Documentation   string             `json:"documentation"`
	TermsOfUse      string             `json:"terms_of_use"`
	TimeLastUpdate  int64              `json:"time_last_update_unix"`
	TimeNextUpdate  int64              `json:"time_next_update_unix"`
	BaseCode        string             `json:"base_code"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

func (c *RateCache) RefreshRates() error {
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", c.apiKey, c.base)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data ExchangeRateAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	if data.Result != "success" {
		return errors.New("failed to fetch rates from API")
	}

	// Filter only supported currencies for caching
	filteredRates := make(map[string]float64)
	for cur := range supported {
		if rate, ok := data.ConversionRates[cur]; ok {
			filteredRates[cur] = rate
		}
	}

	c.mu.Lock()
	c.rates = filteredRates
	c.lastFetch = time.Now()
	c.mu.Unlock()

	return nil
}

func (c *RateCache) GetRate(from, to string) (float64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if from == to {
		return 1, nil
	}

	// Rate(from->to) = Rate(base->to) / Rate(base->from)
	rateFrom, okFrom := c.rates[from]
	rateTo, okTo := c.rates[to]

	if !okFrom || !okTo {
		return 0, errors.New("rate not found")
	}

	return rateTo / rateFrom, nil
}

func validateDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Now(), nil
	}
	const layout = "2006-01-02"
	d, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, use YYYY-MM-DD")
	}
	if d.After(time.Now()) {
		return time.Time{}, errors.New("date cannot be in the future")
	}
	if time.Since(d).Hours() > 24*90 {
		return time.Time{}, errors.New("date older than 90 days not supported")
	}
	return d, nil
}

func convertHandler(cache *RateCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		amountStr := r.URL.Query().Get("amount")
		dateStr := r.URL.Query().Get("date")

		if !supported[from] || !supported[to] {
			http.Error(w, "unsupported currency", http.StatusBadRequest)
			return
		}

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount <= 0 {
			http.Error(w, "invalid amount", http.StatusBadRequest)
			return
		}

		_, err = validateDate(dateStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rate, err := cache.GetRate(from, to)
		if err != nil {
			http.Error(w, "exchange rate not found", http.StatusNotFound)
			return
		}

		result := amount * rate

		resp := map[string]float64{"amount": result}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func main() {
	apiKey := "a123fc8630d96a14097c1594"
	baseCurrency := "USD"

	cache := NewRateCache(apiKey, baseCurrency)

	if err := cache.RefreshRates(); err != nil {
		log.Fatalf("Failed to refresh rates: %v", err)
	}

	// // Periodically refresh every hour
	// go func() {
	// 	for {
	// 		time.Sleep(time.Hour)
	// 		err := cache.RefreshRates()
	// 		if err != nil {
	// 			log.Printf("Error refreshing rates: %v", err)
	// 		} else {
	// 			log.Println("Exchange rates updated.")
	// 		}
	// 	}
	// }()

	http.HandleFunc("/convert", convertHandler(cache))
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
