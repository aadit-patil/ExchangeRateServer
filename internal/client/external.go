package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aadit-patil/ExchangeRateServer/internal/metrics"
)

type ExchangeRateAPIResponse struct {
	ConversionRates map[string]float64 `json:"rates"`
}

func FetchRatesForBase(base, date string) (map[string]float64, error) {

	url := fmt.Sprintf("https://api.frankfurter.dev/v1/%s?base=%s", date, base)
	metrics.APIRequests.Inc()
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data ExchangeRateAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.ConversionRates, err
}
