package api

import (
	"encoding/json"
	"strconv"

	"net/http"

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
	if amount == "" {
		rate, err := service.ConvertCurrency(from, to, date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(model.ConvertResponse{Rate: rate})
	}
	//fetch amt
	amt, err2 := strconv.ParseFloat(amount, 64)
	if err2 != nil {
		json.NewEncoder(w).Encode(model.ConvertResponse{Amount: 0})
		return
	}

	convertedAmt, rate, err := service.ConvertCurrencyWithAmount(amt, from, to, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(model.ConvertResponse{Rate: rate, Amount: convertedAmt})

}
