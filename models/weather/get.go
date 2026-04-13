package weather

import (
	"errors"
	"github.com/shopspring/decimal"
)

type CurrentWeather struct {
	Temperature decimal.Decimal `json:"temperature"`
}

type DailyForecast struct {
	Date        string          `json:"date"`
	Temperature decimal.Decimal `json:"temperature"`
	MinTemp     decimal.Decimal `json:"min_temperature"`
	MaxTemp     decimal.Decimal `json:"max_temperature"`
	Description string          `json:"description"`
}

type Location struct {
	Lat decimal.Decimal `json:"lat"`
	Lon decimal.Decimal `json:"lon"`
}

type LocationWeather struct {
	Lat         decimal.Decimal `json:"lat"`
	Lon         decimal.Decimal `json:"lon"`
	Temperature decimal.Decimal `json:"temperature"`
	Error       string          `json:"error,omitempty"`
}

type Geocoder struct {
	cities map[string]struct {
		Lat decimal.Decimal
		Lon decimal.Decimal
	}
}

func NewGeocoder() *Geocoder {
	g := &Geocoder{
		cities: make(map[string]struct {
			Lat decimal.Decimal
			Lon decimal.Decimal
		}),
	}
	
	g.cities["Минск"] = struct {
		Lat decimal.Decimal
		Lon decimal.Decimal
	}{Lat: decimal.NewFromFloat(53.9), Lon: decimal.NewFromFloat(27.5667)}
	
	g.cities["Лондон"] = struct {
		Lat decimal.Decimal
		Lon decimal.Decimal
	}{Lat: decimal.NewFromFloat(51.5074), Lon: decimal.NewFromFloat(-0.1278)}
	
	g.cities["Токио"] = struct {
		Lat decimal.Decimal
		Lon decimal.Decimal
	}{Lat: decimal.NewFromFloat(35.6895), Lon: decimal.NewFromFloat(139.6917)}
	
	g.cities["Шанхай"] = struct {
		Lat decimal.Decimal
		Lon decimal.Decimal
	}{Lat: decimal.NewFromFloat(31.2304), Lon: decimal.NewFromFloat(121.4737)}
	
	g.cities["Варшава"] = struct {
		Lat decimal.Decimal
		Lon decimal.Decimal
	}{Lat: decimal.NewFromFloat(52.2297), Lon: decimal.NewFromFloat(21.0122)}
	
	return g
}

func (g *Geocoder) GetCoordinates(cityName string) (decimal.Decimal, decimal.Decimal, error) {
	city, exists := g.cities[cityName]
	if !exists {
		return decimal.Zero, decimal.Zero, errors.New("city not found: " + cityName)
	}
	return city.Lat, city.Lon, nil
}

