package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Тест 1: GET /weather/forecast с provider=openweather возвращает 200
func TestWeatherHandler_1_ForecastOpenWeatherReturns200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	req := httptest.NewRequest("GET", "/api/v1/weather/forecast?lat=53.9&lon=27.5667&provider=openweather", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetForecast(c)

	if w.Code != http.StatusOK {
		t.Errorf("Тест 1 НЕ ПРОЙДЕН: ожидался 200, получен %d", w.Code)
	} else {
		t.Log("Тест 1 ПРОЙДЕН")
	}
}

// Тест 2: GET /weather/forecast с provider=openmeteo возвращает 200
func TestWeatherHandler_2_ForecastOpenMeteoReturns200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	req := httptest.NewRequest("GET", "/api/v1/weather/forecast?lat=53.9&lon=27.5667&provider=openmeteo", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetForecast(c)

	if w.Code != http.StatusOK {
		t.Errorf("Тест 2 НЕ ПРОЙДЕН: ожидался 200, получен %d", w.Code)
	} else {
		t.Log("Тест 2 ПРОЙДЕН")
	}
}

// Тест 3: GET /weather/forecast без provider возвращает 400
func TestWeatherHandler_3_ForecastMissingProviderReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	req := httptest.NewRequest("GET", "/api/v1/weather/forecast?lat=53.9&lon=27.5667", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetForecast(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Тест 3 НЕ ПРОЙДЕН: ожидался 400, получен %d", w.Code)
	} else {
		t.Log("Тест 3 ПРОЙДЕН")
	}
}

// Тест 4: GET /weather/forecast с неверным provider возвращает 400
func TestWeatherHandler_4_ForecastInvalidProviderReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	req := httptest.NewRequest("GET", "/api/v1/weather/forecast?lat=53.9&lon=27.5667&provider=invalid", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetForecast(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Тест 4 НЕ ПРОЙДЕН: ожидался 400, получен %d", w.Code)
	} else {
		t.Log("Тест 4 ПРОЙДЕН")
	}
}

// Тест 5: GET /weather/forecast без координат возвращает 400
func TestWeatherHandler_5_ForecastMissingCoordinatesReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	req := httptest.NewRequest("GET", "/api/v1/weather/forecast?provider=openweather", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetForecast(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Тест 5 НЕ ПРОЙДЕН: ожидался 400, получен %d", w.Code)
	} else {
		t.Log("Тест 5 ПРОЙДЕН")
	}
}