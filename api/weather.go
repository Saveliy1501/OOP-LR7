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

type WeatherHandler struct{}

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

// HandleGetForecast godoc
// @Summary      Get Weather Forecast
// @Description  Returns weather forecast for 5 days
// @Tags         weather
// @Produce      json
// @Param        lat      query     string  true  "Latitude" default(53.9)
// @Param        lon      query     string  true  "Longitude" default(27.5667)
// @Param        provider query     string  true  "Weather provider" Enums(openweather, openmeteo)
// @Success      200      {object}  responses.SuccessResponse[[]weather.DailyForecast]
// @Failure      400      {object}  responses.StatusResponse
// @Failure      500      {object}  responses.StatusResponse
// @Router       /weather/forecast [get]
func (h *WeatherHandler) HandleGetForecast(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	providerParam := c.Query("provider")

	if latStr == "" || lonStr == "" {
		c.JSON(400, responses.StatusResponse{Code: 400, Message: "missing coordinates"})
		return
	}

	if providerParam == "" {
		c.JSON(400, responses.StatusResponse{Code: 400, Message: "missing provider parameter"})
		return
	}

	lat, errLat := decimal.NewFromString(latStr)
	lon, errLon := decimal.NewFromString(lonStr)

	if errLat != nil || errLon != nil {
		c.JSON(400, responses.StatusResponse{Code: 400, Message: "invalid coordinates"})
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
			"https://api.openweathermap.org",
		)
	default:
		c.JSON(400, responses.StatusResponse{Code: 400, Message: "unknown provider: " + providerParam})
		return
	}

	controller := controllers.NewCurrentWeatherController(weatherClient)
	forecast, err := controller.GetForecast(lat, lon)

	if err != nil {
		c.JSON(500, responses.StatusResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(200, responses.SuccessResponse[[]weather.DailyForecast]{
		Code:    200,
		Message: "Success",
		Data:    forecast,
	})
}

// HandleGetMultipleLocations 
// @Summary      Get Weather for Multiple Locations
// @Description  Returns current temperature for multiple coordinates
// @Tags         weather
// @Accept       json
// @Produce      json
// @Param        locations body     []weather.Location  true  "Array of locations"
// @Param        provider  query    string               true  "Weather provider" Enums(openweather, openmeteo)
// @Success      200       {object}  responses.SuccessResponse[[]weather.LocationWeather]
// @Failure      400       {object}  responses.StatusResponse
// @Router       /weather/multiple [post]
func (h *WeatherHandler) HandleGetMultipleLocations(c *gin.Context) {
	// ЗАГЛУШКА: возвращаем 501 Not Implemented
	c.JSON(501, responses.StatusResponse{Code: 501, Message: "not implemented"})
}

