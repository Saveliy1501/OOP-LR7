package clients

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
)


// Тест 1: Успешный ответ от сервера
func TestOpenWeatherClient_1_SuccessfulResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"main":{"temp":22.5}}`))
	}))
	defer mockServer.Close()

	client := NewOpenWeatherClient("test_key", mockServer.URL)
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
func TestOpenWeatherClient_2_NegativeTemperature(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"main":{"temp":-15.5}}`))
	}))
	defer mockServer.Close()

	client := NewOpenWeatherClient("test_key", mockServer.URL)
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

// Тест 3: Ошибка авторизации (401)
func TestOpenWeatherClient_3_UnauthorizedError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Invalid API key"}`))
	}))
	defer mockServer.Close()

	client := NewOpenWeatherClient("wrong_key", mockServer.URL)
	lat := decimal.NewFromFloat(53.9)
	lon := decimal.NewFromFloat(27.5667)

	temp, err := client.LocationCurrentTemperature(lat, lon)

	if err == nil {
		t.Error("Тест 3 НЕ ПРОЙДЕН: ожидалась ошибка авторизации, получена nil")
	}
	if !temp.IsZero() {
		t.Errorf("Тест 3 НЕ ПРОЙДЕН: ожидался ноль, получено %v", temp)
	} else {
		t.Log("Тест 3 ПРОЙДЕН: ошибка авторизации обработана")
	}
}

// Тест 4: Пустой ответ от сервера
func TestOpenWeatherClient_4_EmptyResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(``))
	}))
	defer mockServer.Close()

	client := NewOpenWeatherClient("test_key", mockServer.URL)
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
func TestOpenWeatherClient_5_MissingTemperatureField(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"main":{}}`))
	}))
	defer mockServer.Close()

	client := NewOpenWeatherClient("test_key", mockServer.URL)
	lat := decimal.NewFromFloat(53.9)
	lon := decimal.NewFromFloat(27.5667)

	temp, err := client.LocationCurrentTemperature(lat, lon)

	if err == nil {
		t.Logf("Вернуло без ошибки, получено %v", err)
	}
	if temp.IsZero() {
		t.Log("Тест 5 ПРОЙДЕН: отсутствие поля temperature обработано")
	} else {
		t.Errorf("Тест 5 НЕ ПРОЙДЕН: ожидался ноль, получено %v", temp)
	}
}

// Тест 6: Сохранение точности координат
func TestOpenWeatherClient_6_PreservesCoordinatePrecision(t *testing.T) {
	var capturedLat, capturedLon string

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedLat = r.URL.Query().Get("lat")
		capturedLon = r.URL.Query().Get("lon")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"main":{"temp":22.5}}`))
	}))
	defer mockServer.Close()

	client := NewOpenWeatherClient("test_key", mockServer.URL)
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
func TestOpenWeatherClient_7_ValidURLConstruction(t *testing.T) {
    var capturedURL string

    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        capturedURL = r.URL.String()
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"main":{"temp":0}}`))
    }))
    defer mockServer.Close()

    client := NewOpenWeatherClient("abc123key", mockServer.URL)
    lat := decimal.NewFromFloat(53.9)
    lon := decimal.NewFromFloat(27.5667)

    client.LocationCurrentTemperature(lat, lon)

    expected := "/?lat=53.9&lon=27.5667&appid=abc123key&units=metric"
    if capturedURL != expected {
        t.Errorf("Тест 7 НЕ ПРОЙДЕН: ожидался URL %s, получен %s", expected, capturedURL)
    } else {
        t.Log("Тест 7 ПРОЙДЕН: URL сформирован правильно")
    }
}