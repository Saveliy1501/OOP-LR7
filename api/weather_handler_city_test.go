package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Тест 1: Валидный город Минск возвращает 200
func TestWeatherHandler_1_ValidCityMinsk(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	req := httptest.NewRequest("GET", "/api/v1/weather/city?city=Минск&provider=openmeteo", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetWeatherByCity(c)

	if w.Code != http.StatusOK {
		t.Errorf("Тест 1 НЕ ПРОЙДЕН: ожидался 200, получен %d", w.Code)
	} else {
		t.Log("Тест 1 ПРОЙДЕН")
	}
}

// Тест 2: Неизвестный город возвращает 404
func TestWeatherHandler_2_InvalidCity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	req := httptest.NewRequest("GET", "/api/v1/weather/city?city=НесуществующийГород&provider=openmeteo", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetWeatherByCity(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("Тест 2 НЕ ПРОЙДЕН: ожидался 404, получен %d", w.Code)
	} else {
		t.Log("Тест 2 ПРОЙДЕН")
	}
}

// Тест 3: Отсутствует параметр city возвращает 400
func TestWeatherHandler_3_MissingCity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	req := httptest.NewRequest("GET", "/api/v1/weather/city?provider=openmeteo", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetWeatherByCity(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Тест 3 НЕ ПРОЙДЕН: ожидался 400, получен %d", w.Code)
	} else {
		t.Log("Тест 3 ПРОЙДЕН")
	}
}