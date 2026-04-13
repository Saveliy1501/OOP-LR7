 package controllers

import (
	"clients"
	weather "models/weather"

	"github.com/shopspring/decimal"
)

type CurrentWeatherController[T clients.WeatherDataClient] struct {
	Client T
}

func NewCurrentWeatherController[T clients.WeatherDataClient](client T) *CurrentWeatherController[T] {
	return &CurrentWeatherController[T]{
		Client: client,
	}
}

func (c *CurrentWeatherController[T]) GetCurrentWeather(lat decimal.Decimal, lon decimal.Decimal) (weather.CurrentWeather, error) {

	temperature, err := c.Client.LocationCurrentTemperature(lat, lon)
	if err != nil {
		return weather.CurrentWeather{}, err
	}

	return weather.CurrentWeather{
		Temperature: temperature,
	}, nil
}

func (c *CurrentWeatherController[T]) GetForecast(lat decimal.Decimal, lon decimal.Decimal) ([]weather.DailyForecast, error) {
	forecast, err := c.Client.LocationForecast(lat, lon)
	if err != nil {
		return nil, err
	}
	return forecast, nil
}

func (c *CurrentWeatherController[T]) GetMultipleCurrentWeather(locations []weather.Location) []weather.LocationWeather {
	// Заглушка для RED стадии
	return []weather.LocationWeather{}
}