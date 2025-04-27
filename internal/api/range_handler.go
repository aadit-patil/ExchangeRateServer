package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aadit-patil/ExchangeRateServer/internal/errors"
	"github.com/aadit-patil/ExchangeRateServer/internal/service"
	"github.com/gammazero/workerpool"
)

type RangeResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func ConvertRangeHandler(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	amount := r.URL.Query().Get("amount")

	if from == "" || to == "" || start == "" || end == "" {
		errors.ErrBadRequest.Write(w)
		return
	}

	startDate, err1 := time.Parse("2006-01-02", start)
	endDate, err2 := time.Parse("2006-01-02", end)
	if err1 != nil || err2 != nil || startDate.After(endDate) {
		errors.ErrInvalidDateRange.Write(w)
		return
	}

	if !service.IsValidDate(start) || !service.IsValidDate(end) {
		errors.ErrInvalidDateRange.Write(w)
		return
	}

	var multiplier float64 = 1.0
	if amount != "" {
		m, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			errors.ErrInvalidAmount.Write(w)
			return
		}
		multiplier = m
	}

	rates := make(map[string]float64)
	//workerPools with configured workes = number of days
	wp := workerpool.New(service.DaysBetween(startDate, endDate))
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")

		wp.Submit(func() {
			service.GetGlobalStrategy().FetchAndStoreRate(from, to, dateStr)
		})

	}

	wp.StopWait()

	//fetch from cache
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		rate, err := service.ConvertCurrency(from, to, dateStr)
		if err == nil {
			rates[dateStr] = rate * multiplier
		}
	}

	json.NewEncoder(w).Encode(RangeResponse{Rates: rates})
}
