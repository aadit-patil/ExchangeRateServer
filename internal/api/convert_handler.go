package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aadit-patil/ExchangeRateServer/internal/configs"
	"github.com/aadit-patil/ExchangeRateServer/internal/errors"
	"github.com/aadit-patil/ExchangeRateServer/internal/service"
)

type ConvertResponse struct {
	Rate   float64 `json:"rate"`
	Amount float64 `json:"amount,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	date := r.URL.Query().Get("date")
	amount := r.URL.Query().Get("amount")

	if from == "" || to == "" {
		errors.ErrBadRequest.Write(w)
		return
	}

	if configs.SupportedCurrencies[from] || configs.SupportedCurrencies[to] {
		errors.ErrUnsupportedCurrency.Write(w)
		return
	}

	if !service.IsValidDate(date) {
		errors.ErrInvalidDateRange.Write(w)
		return
	}

	amt, err2 := strconv.ParseFloat(amount, 64)
	if err2 != nil {
		errors.ErrInvalidAmount.Write(w)
		return
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	rate, err := service.ConvertCurrency(from, to, date)
	if err != nil {
		// start async background update
		go service.GetGlobalStrategy().FetchAndStoreRate(from, to, date)
		//use prefetched rate
		rate, _ = service.ConvertCurrency(from, to, time.Now().Format("2006-01-02"))
	}

	if amount == "" {
		json.NewEncoder(w).Encode(ConvertResponse{Rate: rate})
		return
	}
	convertedAmt := amt * rate
	json.NewEncoder(w).Encode(ConvertResponse{Rate: rate, Amount: convertedAmt})

}
