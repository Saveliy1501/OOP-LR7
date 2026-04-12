package clients

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
)

// Тест 1: OpenWeather прогноз возвращает 5 дней
func TestOpenWeatherClient_1_ForecastReturns5Days(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list": [
				{"dt_txt": "2025-04-13 12:00:00", "main": {"temp": 15.5, "temp_min": 10.0, "temp_max": 18.0}, "weather": [{"description": "clear"}]},
				{"dt_txt": "2025-04-14 12:00:00", "main": {"temp": 14.0, "temp_min": 9.0, "temp_max": 16.5}, "weather": [{"description": "rain"}]},
				{"dt_txt": "2025-04-15 12:00:00", "main": {"temp": 12.5, "temp_min": 7.5, "temp_max": 14.0}, "weather": [{"description": "clouds"}]},
				{"dt_txt": "2025-04-16 12:00:00", "main": {"temp": 13.0, "temp_min": 8.0, "temp_max": 15.0}, "weather": [{"description": "clouds"}]},
				{"dt_txt": "2025-04-17 12:00:00", "main": {"temp": 11.0, "temp_min": 6.0, "temp_max": 13.0}, "weather": [{"description": "rain"}]}
			]
		}`))
	}))
	defer mockServer.Close()

	client := NewOpenWeatherClient("test_key", mockServer.URL)
	lat := decimal.NewFromFloat(53.9)
	lon := decimal.NewFromFloat(27.5667)

	forecast, err := client.LocationForecast(lat, lon)

	if err != nil {
		t.Errorf("Тест 1 НЕ ПРОЙДЕН: ошибка %v", err)
	}
	if len(forecast) != 5 {
		t.Errorf("Тест 1 НЕ ПРОЙДЕН: ожидалось 5 дней, получено %d", len(forecast))
	} else {
		t.Log("Тест 1 ПРОЙДЕН")
	}
}

// Тест 2: OpenMeteo прогноз возвращает 5 дней
func TestOpenMeteoClient_2_ForecastReturns5Days(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"daily": {
				"time": ["2025-04-13", "2025-04-14", "2025-04-15", "2025-04-16", "2025-04-17"],
				"temperature_2m_max": [15.5, 14.0, 12.5, 13.0, 11.0],
				"temperature_2m_min": [10.0, 9.0, 7.5, 8.0, 6.0]
			}
		}`))
	}))
	defer mockServer.Close()

	client := NewOpenMeteoClient(mockServer.URL)
	lat := decimal.NewFromFloat(53.9)
	lon := decimal.NewFromFloat(27.5667)

	forecast, err := client.LocationForecast(lat, lon)

	if err != nil {
		t.Errorf("Тест 2 НЕ ПРОЙДЕН: ошибка %v", err)
	}
	if len(forecast) != 5 {
		t.Errorf("Тест 2 НЕ ПРОЙДЕН: ожидалось 5 дней, получено %d", len(forecast))
	} else {
		t.Log("Тест 2 ПРОЙДЕН")
	}
}

// Тест 3: OpenWeather правильно формирует URL для прогноза
func TestOpenWeatherClient_3_ForecastURLConstruction(t *testing.T) {
	var capturedURL string

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"list":[]}`))
	}))
	defer mockServer.Close()

	client := NewOpenWeatherClient("abc123key", mockServer.URL)
	lat := decimal.NewFromFloat(53.9)
	lon := decimal.NewFromFloat(27.5667)

	client.LocationForecast(lat, lon)

	expected := "/data/2.5/forecast?appid=abc123key&lat=53.9&lon=27.5667&units=metric&cnt=40"
	if capturedURL != expected {
		t.Errorf("Тест 3 НЕ ПРОЙДЕН: ожидался URL %s, получен %s", expected, capturedURL)
	} else {
		t.Log("Тест 3 ПРОЙДЕН: URL для прогноза OpenWeather сформирован правильно")
	}
}

// Тест 4: OpenMeteo правильно формирует URL для прогноза
func TestOpenMeteoClient_4_ForecastURLConstruction(t *testing.T) {
	var capturedURL string

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"daily":{"time":[],"temperature_2m_max":[],"temperature_2m_min":[]}}`))
	}))
	defer mockServer.Close()

	client := NewOpenMeteoClient(mockServer.URL)
	lat := decimal.NewFromFloat(53.9)
	lon := decimal.NewFromFloat(27.5667)

	client.LocationForecast(lat, lon)

	expected := "/?latitude=53.9&longitude=27.5667&daily=temperature_2m_max,temperature_2m_min&timezone=auto&days=5"
	if capturedURL != expected {
		t.Errorf("Тест 4 НЕ ПРОЙДЕН: ожидался URL %s, получен %s", expected, capturedURL)
	} else {
		t.Log("Тест 4 ПРОЙДЕН: URL для прогноза OpenMeteo сформирован правильно")
	}
}