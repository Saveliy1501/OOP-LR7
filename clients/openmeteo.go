package clients

import (
	"encoding/json"
	"fmt"
	"net/http"

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
