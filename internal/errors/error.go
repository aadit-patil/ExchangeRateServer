package errors

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func (e *AppError) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	_ = json.NewEncoder(w).Encode(e)
}

var (
	ErrBadRequest          = New(http.StatusBadRequest, "Missing required parameters")
	ErrInvalidAmount       = New(http.StatusBadRequest, "Invalid amount")
	ErrInternal            = New(http.StatusInternalServerError, "Internal server error")
	ErrInvalidDateRange    = New(http.StatusBadRequest, "Date must be within the past 90 days")
	ErrNotFound            = New(http.StatusNotFound, "Data not found")
	ErrUnsupportedCurrency = New(http.StatusNotFound, "Unsupported Currency")
)
