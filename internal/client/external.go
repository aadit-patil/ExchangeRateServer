package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aadit-patil/ExchangeRateServer/internal/metrics"
)

type ExchangeRateAPIResponse struct {
	ConversionRates    map[string]float64 `json:"conversion_rates"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
}

func FetchRatesForBase(base string) (map[string]float64, time.Time, error) {
	apiKey := "a123fc8630d96a14097c1594"
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", apiKey, base)
	metrics.APIRequests.Inc()
	resp, err := http.Get(url)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer resp.Body.Close()

	var data ExchangeRateAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	expiresAt := time.Unix(data.TimeNextUpdateUnix, 0)
	return data.ConversionRates, expiresAt, err
}
