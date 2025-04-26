package service

var supportedCurrencies = map[string]bool{
	"USD": true,
	"INR": true,
	"EUR": true,
	"JPY": true,
	"GBP": true,
}

func ConvertCurrency(from, to, date string) (float64, error) {
	return strategy.GetRate(from, to, date)
}

func ConvertCurrencyWithAmount(amount float64, from, to, date string) (float64, float64, error) {
	rate, err := strategy.GetRate(from, to, date)
	if err != nil {
		return 0, 0, err
	}
	return rate * amount, rate, nil
}
