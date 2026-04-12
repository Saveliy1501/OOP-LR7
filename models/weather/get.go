package weather

import "github.com/shopspring/decimal"

type CurrentWeather struct {
	Temperature decimal.Decimal `json:"temperature"`
}

type DailyForecast struct {
	Date        string          `json:"date"`
	Temperature decimal.Decimal `json:"temperature"`
	MinTemp     decimal.Decimal `json:"min_temperature"`
	MaxTemp     decimal.Decimal `json:"max_temperature"`
	Description string          `json:"description"`
}

