package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Тест 1: провайдер openmeteo возвращает 200
func TestWeatherHandler_1_ProviderOpenMeteoReturns200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	req := httptest.NewRequest("GET", "/api/v1/weather?lat=53.9&lon=27.5667&provider=openmeteo", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetCurrentWeather(c)

	if w.Code != http.StatusOK {
		t.Errorf("Тест 1 НЕ ПРОЙДЕН: ожидался 200, получен %d", w.Code)
	} else {
		t.Log("Тест 1 ПРОЙДЕН")
	}
}

// Тест 2: неверный провайдер возвращает 400
func TestWeatherHandler_2_InvalidProviderReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	req := httptest.NewRequest("GET", "/api/v1/weather?lat=53.9&lon=27.5667&provider=invalid", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetCurrentWeather(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Тест 2 НЕ ПРОЙДЕН: ожидался 400, получен %d", w.Code)
	} else {
		t.Log("Тест 2 ПРОЙДЕН")
	}
}

// Тест 3: регистр провайдера не важен
func TestWeatherHandler_3_ProviderCaseInsensitive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	req := httptest.NewRequest("GET", "/api/v1/weather?lat=53.9&lon=27.5667&provider=OpenMeteo", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetCurrentWeather(c)

	if w.Code != http.StatusOK {
		t.Errorf("Тест 3 НЕ ПРОЙДЕН: ожидался 200, получен %d", w.Code)
	} else {
		t.Log("Тест 3 ПРОЙДЕН")
	}
}