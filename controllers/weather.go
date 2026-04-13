 package controllers

import (
	"clients"
	weather "models/weather"
	"sync"
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
	results := make([]weather.LocationWeather, len(locations))
	var wg sync.WaitGroup

	for i, loc := range locations {
		wg.Add(1)
		go func(idx int, location weather.Location) {
			defer wg.Done()

			results[idx].Lat = location.Lat
			results[idx].Lon = location.Lon

			temp, err := c.Client.LocationCurrentTemperature(location.Lat, location.Lon)
			if err != nil {
				results[idx].Error = err.Error()
			} else {
				results[idx].Temperature = temp
			}
		}(i, loc)
	}

	wg.Wait()
	return results
}