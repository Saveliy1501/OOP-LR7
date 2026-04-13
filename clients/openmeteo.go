package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"models/weather"

	"github.com/shopspring/decimal"
)

type openMeteoResponse struct {
	CurrentWeather struct {
		Temperature decimal.Decimal `json:"temperature"`
	} `json:"current_weather"`
}

type OpenMeteoClient struct {
	baseURL string
}

func NewOpenMeteoClient(baseURL string) *OpenMeteoClient {
	if baseURL == "" {
		baseURL = "https://api.open-meteo.com/v1/forecast"
	}
	return &OpenMeteoClient{
		baseURL: baseURL,
	}
}

func (c *OpenMeteoClient) LocationCurrentTemperature(lat decimal.Decimal, lon decimal.Decimal) (decimal.Decimal, error) {
	url := fmt.Sprintf("%s?latitude=%s&longitude=%s&current_weather=true",
		c.baseURL, lat.String(), lon.String())

	resp, err := http.Get(url)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to call open-meteo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("open-meteo returned bad status: %d", resp.StatusCode)
	}

	var data openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return decimal.Zero, fmt.Errorf("failed to decode response: %w", err)
	}

	return data.CurrentWeather.Temperature, nil
}

func (c *OpenMeteoClient) LocationForecast(lat decimal.Decimal, lon decimal.Decimal) ([]weather.DailyForecast, error) {
	url := fmt.Sprintf("%s?latitude=%s&longitude=%s&daily=temperature_2m_max,temperature_2m_min&timezone=auto&days=5",
		c.baseURL, lat.String(), lon.String())

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call openmeteo forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openmeteo forecast returned bad status: %d", resp.StatusCode)
	}

	var data struct {
		Daily struct {
			Time             []string          `json:"time"`
			Temperature2mMax []decimal.Decimal `json:"temperature_2m_max"`
			Temperature2mMin []decimal.Decimal `json:"temperature_2m_min"`
		} `json:"daily"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	var forecast []weather.DailyForecast
	for i := 0; i < len(data.Daily.Time); i++ {
		sum := data.Daily.Temperature2mMax[i].Add(data.Daily.Temperature2mMin[i])
		avgTemp := sum.Div(decimal.NewFromFloat(2))

		forecast = append(forecast, weather.DailyForecast{
			Date:        data.Daily.Time[i],
			Temperature: avgTemp,
			MinTemp:     data.Daily.Temperature2mMin[i],
			MaxTemp:     data.Daily.Temperature2mMax[i],
			Description: "",
		})
	}

	return forecast, nil
}
