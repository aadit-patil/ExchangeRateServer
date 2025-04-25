package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// can be used for future development and validations
type ExchangeRateAPIResponse struct {
	Result          string             `json:"result"`
	Documentation   string             `json:"documentation"`
	TermsOfUse      string             `json:"terms_of_use"`
	TimeLastUpdate  int64              `json:"time_last_update_unix"`
	TimeNextUpdate  int64              `json:"time_next_update_unix"`
	BaseCode        string             `json:"base_code"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

func FetchRate(from, to, date string) (float64, error) {
	apiKey := "a123fc8630d96a14097c1594"
	base := "USD"
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", apiKey, base)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data ExchangeRateAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.ConversionRates[to], err
}
