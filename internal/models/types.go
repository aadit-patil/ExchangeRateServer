package model

type ConvertRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Date   string `json:"date"`
	Amount string `json:"amount"`
}

type ConvertResponse struct {
	Rate   float64 `json:"rate"`
	Amount float64 `json:"amount,omitempty"`
}
