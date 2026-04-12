package clients

import (
	"github.com/shopspring/decimal"
	"models/weather"
)

type WeatherDataClient interface {
	LocationCurrentTemperature(lat decimal.Decimal, lon decimal.Decimal) ( decimal.Decimal,  error)
	LocationForecast(lat decimal.Decimal, lon decimal.Decimal) ([]weather.DailyForecast, error)
}
