package model

type ConvertRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
	Date string `json:"date"`
}

type ConvertResponse struct {
	Rate float64 `json:"rate"`
}
