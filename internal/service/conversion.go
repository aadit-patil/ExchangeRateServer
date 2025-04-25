package service

import "fmt"

func ConvertCurrency(from, to, date string) (float64, error) {
	return strategy.GetRate(from, to, date)
}

func ConvertCurrencyWithAmount(amount float64, from, to, date string) (float64, float64, error) {
	rate, err := strategy.GetRate(from, to, date)
	if err != nil {
		fmt.Println("Error While Fetching Rate:", err.Error())
		return 0, 0, err
	}
	return rate * amount, rate, nil

}
