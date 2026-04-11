package api

import (
	"clients"
	"controllers"
	weather "models/weather"
	"shared/responses"
	"shared/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)
 
type WeatherHandler struct {}
	

func NewCurrentWeatherHandler() *WeatherHandler {
	return &WeatherHandler{}
}

// GetCurrentWeather godoc
// @Summary      Get Current Weather
// @Description  Returns current weather for given coordinates
// @Tags         weather
// @Produce      json
// @Param        lat   query     string  true  "Latitude"    default(18.300231990440125)
// @Param        lon   query     string  true  "Longitude"   default(-64.8251590359234)
// @Param        provider query     string  true  "Weather provider" Enums(openweather, openmeteo)
// @Success      200   {object}  responses.SuccessResponse[weather.CurrentWeather]
// @Failure      400   {object}  responses.StatusResponse
// @Failure      500   {object}  responses.StatusResponse
// @Router       /weather [get]
func (h *WeatherHandler) HandleGetCurrentWeather(c *gin.Context) {
	lat, errLat := decimal.NewFromString(c.Query("lat"))
	lon, errLon := decimal.NewFromString(c.Query("lon"))

	if errLat != nil || errLon != nil {
		c.JSON(400, responses.StatusResponse{Code: 400, Message: "invalid coordinates"})
		return
	}

	providerParam := c.Query("provider")
	if providerParam == "" {
		c.JSON(400, responses.StatusResponse{Code: 400, Message: "missing provider parameter"})
		return
	}

	providerParam = strings.ToLower(providerParam)

	var weatherClient clients.WeatherDataClient

	switch providerParam {
	case "openmeteo":
		weatherClient = clients.NewOpenMeteoClient(
			utils.GetEnv("OPENMETEO_BASE_URL", "https://api.open-meteo.com/v1/forecast"),
		)
	case "openweather":
		weatherClient = clients.NewOpenWeatherClient(
			utils.GetEnv("OPENWEATHER_API_KEY", ""),
			utils.GetEnv("OPENWEATHER_BASE_URL", "https://api.openweathermap.org/data/2.5/weather"),
		)
	default:
		c.JSON(400, responses.StatusResponse{Code: 400, Message: "unknown provider: " + providerParam})
		return
	}

	controller := controllers.NewCurrentWeatherController(weatherClient)
	result, err := controller.GetCurrentWeather(lat, lon)

	if err != nil {

		c.JSON(500, responses.StatusResponse{Code: 500, Message: err.Error()})
		return
	}
	c.JSON(200, responses.SuccessResponse[weather.CurrentWeather]{Code: 200, Message: "Success", Data: result})
}
