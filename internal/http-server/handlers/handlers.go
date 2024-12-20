package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"weatherapp/internal/api_client"
	"weatherapp/internal/config"
	"weatherapp/internal/models"
	"weatherapp/internal/storage"
)

type Handlers struct {
	Storage *storage.Storage
	Cfg     config.Config
}

func NewWeatherHandler(storage *storage.Storage, cfg config.Config) *Handlers {
	return &Handlers{Storage: storage, Cfg: cfg}
}

func (h *Handlers) GetByCityName(w http.ResponseWriter, r *http.Request) {

	var req struct {
		CityName string `json:"city_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CityName == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	weather, err := h.Storage.GetWeather(req.CityName)
	if err == nil && len(weather) > 0 {
		if err := json.NewEncoder(w).Encode(weather[0]); err != nil {
			log.Printf("Failed to encode weather data to JSON: %v", err)
		}

		return
	}

	weatherData, err := api_client.GetWeatherByURL(req.CityName, h.Cfg)
	if err != nil {
		log.Printf("Failed to fetch weather data from API: %s", err)
		return
	}

	if err := h.Storage.SaveWeather(weatherData.CityName, weatherData.Temperature, weatherData.FeelsLike, weatherData.Pressure, weatherData.Humidity); err != nil {
		log.Printf("Failed to save weather data: %s", err)
	}

	if err := json.NewEncoder(w).Encode(weatherData); err != nil {
		log.Printf("Failed to encode weather data to JSON: %v", err)
		return
	}
}

func (h *Handlers) GetPopularCities(w http.ResponseWriter, r *http.Request) {

	popularCities := []string{"Moscow", "Irkutsk", "London", "Texas", "Tokyo"}

	var response []models.Weather

	for _, city := range popularCities {
		weather, err := h.Storage.GetWeather(city)
		if err == nil && len(weather) > 0 {
			response = append(response, models.Weather{
				CityName:    weather[0].CityName,
				Temperature: weather[0].Temperature,
				FeelsLike:   weather[0].FeelsLike,
				Pressure:    weather[0].Pressure,
				Humidity:    weather[0].Humidity,
			})
			continue
		}

		weatherData, err := api_client.GetWeatherByURL(city, h.Cfg)
		if err != nil {
			log.Printf("Failed to fetch weather data for city '%s': %s", city, err)
			continue
		}

		if err := h.Storage.SaveWeather(weatherData.CityName, weatherData.Temperature, weatherData.FeelsLike, weatherData.Pressure, weatherData.Humidity); err != nil {
			log.Printf("Failed to save weather data for city '%s': %s", city, err)
		}

		response = append(response, *weatherData)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode weather data to JSON: %v", err)
	}
}
