package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	model "github.com/aadit-patil/ExchangeRateServer/internal/models"
	"github.com/aadit-patil/ExchangeRateServer/internal/service"
)

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	date := r.URL.Query().Get("date")
	amount := r.URL.Query().Get("amount")

	if from == "" || to == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	// If date is empty, use today's date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	if amount == "" {
		rate, err := service.ConvertCurrency(from, to, date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(model.ConvertResponse{Rate: rate})
		return
	}

	amt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	converted, rate, err := service.ConvertCurrencyWithAmount(amt, from, to, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(model.ConvertResponse{Rate: rate, Amount: converted})
}
