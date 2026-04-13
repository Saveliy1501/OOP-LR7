package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"models/weather"
	"strings"

	"github.com/shopspring/decimal"
)

type openWeatherResponse struct {
	Main struct {
		Temp decimal.Decimal `json:"temp"`
	} `json:"main"`
}

type OpenWeatherClient struct {
	apiKey  string
	baseURL string
}

func NewOpenWeatherClient(apiKey string, baseURL string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// Implementation of WeatherDataClient
func (c *OpenWeatherClient) LocationCurrentTemperature(lat decimal.Decimal, lon decimal.Decimal) (decimal.Decimal, error) {
	url := fmt.Sprintf("%s/data/2.5/weather?lat=%s&lon=%s&appid=%s&units=metric",
        c.baseURL, lat.String(), lon.String(), c.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to call openweather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("openweather returned bad status: %d", resp.StatusCode)
	}

	var data openWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return decimal.Zero, fmt.Errorf("failed to decode response: %w", err)
	}

	return data.Main.Temp, nil
}

func (c *OpenWeatherClient) LocationForecast(lat decimal.Decimal, lon decimal.Decimal) ([]weather.DailyForecast, error) {
	url := fmt.Sprintf("%s/data/2.5/forecast?appid=%s&lat=%s&lon=%s&units=metric&cnt=40",
		c.baseURL,c.apiKey, lat.String(), lon.String())
    fmt.Println("DEBUG URL:", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call openweather forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openweather forecast returned bad status: %d", resp.StatusCode)
	}

	var data struct {
		List []struct {
			DtTxt string `json:"dt_txt"`
			Main  struct {
				Temp    decimal.Decimal `json:"temp"`
				TempMin decimal.Decimal `json:"temp_min"`
				TempMax decimal.Decimal `json:"temp_max"`
			} `json:"main"`
			Weather []struct {
				Description string `json:"description"`
			} `json:"weather"`
		} `json:"list"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	var forecast []weather.DailyForecast
dailyData := make(map[string]struct {
    SumTemp     decimal.Decimal
    MinTemp     decimal.Decimal
    MaxTemp     decimal.Decimal
    Description string
    Count       int
})

for _, item := range data.List {
    date := strings.Split(item.DtTxt, " ")[0]
    
    entry := dailyData[date]
    
    entry.SumTemp = entry.SumTemp.Add(item.Main.Temp)
    entry.Count++
    
   
    if entry.Count == 1 || item.Main.TempMin.LessThan(entry.MinTemp) {
        entry.MinTemp = item.Main.TempMin
    }
    
    if entry.Count == 1 || item.Main.TempMax.GreaterThan(entry.MaxTemp) {
        entry.MaxTemp = item.Main.TempMax
    }
    
    if entry.Count == 1 {
        entry.Description = item.Weather[0].Description
    }
    
    dailyData[date] = entry
}

for date, data := range dailyData {
    avgTemp := data.SumTemp.Div(decimal.NewFromInt(int64(data.Count)))
    
    forecast = append(forecast, weather.DailyForecast{
        Date:        date,
        Temperature: avgTemp,
        MinTemp:     data.MinTemp,
        MaxTemp:     data.MaxTemp,
        Description: data.Description,
    })
}

	return forecast, nil
}

