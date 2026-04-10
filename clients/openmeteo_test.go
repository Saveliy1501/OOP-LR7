package clients

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
)

// Тест 1: Успешный ответ от сервера
func TestOpenMeteoClient_1_SuccessfulResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current_weather":{"temperature":22.5}}`))
	}))
	defer mockServer.Close()

	client := NewOpenMeteoClient(mockServer.URL)
	lat := decimal.NewFromFloat(53.9)
	lon := decimal.NewFromFloat(27.5667)

	temp, err := client.LocationCurrentTemperature(lat, lon)

	if err != nil {
		t.Errorf("Тест 1 НЕ ПРОЙДЕН: неожиданная ошибка %v", err)
	}
	expected := decimal.NewFromFloat(22.5)
	if !temp.Equal(expected) {
		t.Errorf("Тест 1 НЕ ПРОЙДЕН: ожидалось %v, получено %v", expected, temp)
	} else {
		t.Log("Тест 1 ПРОЙДЕН: успешный ответ обработан корректно")
	}
}

// Тест 2: Отрицательная температура
func TestOpenMeteoClient_2_NegativeTemperature(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current_weather":{"temperature":-15.5}}`))
	}))
	defer mockServer.Close()

	client := NewOpenMeteoClient(mockServer.URL)
	lat := decimal.NewFromFloat(60.0)
	lon := decimal.NewFromFloat(30.0)

	temp, err := client.LocationCurrentTemperature(lat, lon)

	if err != nil {
		t.Errorf("Тест 2 НЕ ПРОЙДЕН: неожиданная ошибка %v", err)
	}
	expected := decimal.NewFromFloat(-15.5)
	if !temp.Equal(expected) {
		t.Errorf("Тест 2 НЕ ПРОЙДЕН: ожидалось %v, получено %v", expected, temp)
	} else {
		t.Log("Тест 2 ПРОЙДЕН: отрицательная температура обработана")
	}
}

// Тест 3: Ошибка сервера (500)
func TestOpenMeteoClient_3_ServerError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer mockServer.Close()

	client := NewOpenMeteoClient(mockServer.URL)
	lat := decimal.NewFromFloat(53.9)
	lon := decimal.NewFromFloat(27.5667)

	temp, err := client.LocationCurrentTemperature(lat, lon)

	if err == nil {
		t.Error("Тест 3 НЕ ПРОЙДЕН: ожидалась ошибка сервера, получена nil")
	}
	if !temp.IsZero() {
		t.Errorf("Тест 3 НЕ ПРОЙДЕН: ожидался ноль, получено %v", temp)
	} else {
		t.Log("Тест 3 ПРОЙДЕН: ошибка сервера обработана")
	}
}

// Тест 4: Пустой ответ от сервера
func TestOpenMeteoClient_4_EmptyResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(``))
	}))
	defer mockServer.Close()

	client := NewOpenMeteoClient(mockServer.URL)
	lat := decimal.NewFromFloat(53.9)
	lon := decimal.NewFromFloat(27.5667)

	temp, err := client.LocationCurrentTemperature(lat, lon)

	if err == nil {
		t.Error("Тест 4 НЕ ПРОЙДЕН: ожидалась ошибка при пустом ответе, получена nil")
	}
	if !temp.IsZero() {
		t.Errorf("Тест 4 НЕ ПРОЙДЕН: ожидался ноль, получено %v", temp)
	} else {
		t.Log("Тест 4 ПРОЙДЕН: пустой ответ обработан")
	}
}

// Тест 5: Отсутствует поле temperature в ответе
func TestOpenMeteoClient_5_MissingTemperatureField(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current_weather":{}}`))
	}))
	defer mockServer.Close()

	client := NewOpenMeteoClient(mockServer.URL)
	lat := decimal.NewFromFloat(53.9)
	lon := decimal.NewFromFloat(27.5667)

	temp, err := client.LocationCurrentTemperature(lat, lon)

	if err == nil {
		t.Logf("Вернуло без ошибки, получено %v", err)
	}
	if !temp.IsZero() {
		t.Errorf("Тест 5 НЕ ПРОЙДЕН: ожидался ноль, получено %v", temp)
	} else {
		t.Log("Тест 5 ПРОЙДЕН: отсутствие поля temperature обработано")
	}
}

// Тест 6: Сохранение точности координат
func TestOpenMeteoClient_6_PreservesCoordinatePrecision(t *testing.T) {
	var capturedLat, capturedLon string

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedLat = r.URL.Query().Get("latitude")
		capturedLon = r.URL.Query().Get("longitude")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current_weather":{"temperature":22.5}}`))
	}))
	defer mockServer.Close()

	client := NewOpenMeteoClient(mockServer.URL)
	lat := decimal.NewFromFloat(53.904539)
	lon := decimal.NewFromFloat(27.561520)

	client.LocationCurrentTemperature(lat, lon)

	expectedLat := "53.904539"
	expectedLon := "27.56152"

	if capturedLat != expectedLat {
		t.Errorf("Тест 6 НЕ ПРОЙДЕН: ожидалась широта %s, получена %s", expectedLat, capturedLat)
	} else if capturedLon != expectedLon {
		t.Errorf("Тест 6 НЕ ПРОЙДЕН: ожидалась долгота %s, получена %s", expectedLon, capturedLon)
	} else {
		t.Log("Тест 6 ПРОЙДЕН: точность координат сохранена")
	}
}

// Тест 7: Проверка правильности формирования URL
func TestOpenMeteoClient_7_ValidURLConstruction(t *testing.T) {
	var capturedURL string

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current_weather":{"temperature":0}}`))
	}))
	defer mockServer.Close()

	client := NewOpenMeteoClient(mockServer.URL)
	lat := decimal.NewFromFloat(53.9)
	lon := decimal.NewFromFloat(27.5667)

	client.LocationCurrentTemperature(lat, lon)

	expected := "/?latitude=53.9&longitude=27.5667&current_weather=true"
	if capturedURL != expected {
		t.Errorf("Тест 7 НЕ ПРОЙДЕН: ожидался URL %s, получен %s", expected, capturedURL)
	} else {
		t.Log("Тест 7 ПРОЙДЕН: URL сформирован правильно")
	}
}