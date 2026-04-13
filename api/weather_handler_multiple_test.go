package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	
)

func TestWeatherHandler_1_MultipleLocationsReturns200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	locations := []map[string]interface{}{
		{"lat": 53.9, "lon": 27.5667},
		{"lat": 51.5074, "lon": -0.1278},
	}
	jsonBody, _ := json.Marshal(locations)

	req := httptest.NewRequest("POST", "/api/v1/weather/multiple?provider=openmeteo", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetMultipleLocations(c)

	if w.Code != http.StatusOK {
		t.Errorf("Тест 1 НЕ ПРОЙДЕН: ожидался 200, получен %d", w.Code)
	} else {
		t.Log("Тест 1 ПРОЙДЕН")
	}
}

func TestWeatherHandler_2_MultipleLocationsEmptyArray(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	locations := []map[string]interface{}{}
	jsonBody, _ := json.Marshal(locations)

	req := httptest.NewRequest("POST", "/api/v1/weather/multiple?provider=openmeteo", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetMultipleLocations(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Тест 2 НЕ ПРОЙДЕН: ожидался 400, получен %d", w.Code)
	} else {
		t.Log("Тест 2 ПРОЙДЕН")
	}
}

func TestWeatherHandler_3_MultipleLocationsMissingProvider(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	locations := []map[string]interface{}{
		{"lat": 53.9, "lon": 27.5667},
	}
	jsonBody, _ := json.Marshal(locations)

	req := httptest.NewRequest("POST", "/api/v1/weather/multiple", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetMultipleLocations(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Тест 3 НЕ ПРОЙДЕН: ожидался 400, получен %d", w.Code)
	} else {
		t.Log("Тест 3 ПРОЙДЕН")
	}
}

func TestWeatherHandler_4_MultipleLocationsInvalidProvider(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	locations := []map[string]interface{}{
		{"lat": 53.9, "lon": 27.5667},
	}
	jsonBody, _ := json.Marshal(locations)

	req := httptest.NewRequest("POST", "/api/v1/weather/multiple?provider=invalid", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetMultipleLocations(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Тест 4 НЕ ПРОЙДЕН: ожидался 400, получен %d", w.Code)
	} else {
		t.Log("Тест 4 ПРОЙДЕН")
	}
}

func TestWeatherHandler_5_MultipleLocationsInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewCurrentWeatherHandler()

	req := httptest.NewRequest("POST", "/api/v1/weather/multiple?provider=openmeteo", bytes.NewReader([]byte(`{invalid json}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleGetMultipleLocations(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Тест 5 НЕ ПРОЙДЕН: ожидался 400, получен %d", w.Code)
	} else {
		t.Log("Тест 5 ПРОЙДЕН")
	}
}