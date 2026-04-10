package clients

import (
	"github.com/shopspring/decimal"
	
)

type OpenMeteoClient struct {
	baseURL string
}

func NewOpenMeteoClient(baseURL string) *OpenMeteoClient {
	return &OpenMeteoClient{
		baseURL: baseURL,
	}
}

func (c *OpenMeteoClient) LocationCurrentTemperature(lat decimal.Decimal, lon decimal.Decimal) (decimal.Decimal, error) {
	// TODO: реализовать
	return decimal.NewFromFloat(999.99), nil
}
