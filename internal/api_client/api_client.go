package api_client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"weatherapp/internal/config"
	"weatherapp/internal/models"
)

type Response struct {
	Main map[string]float32 `json:"main"`
}

func GetWeatherByURL(cityName string, cfg config.Config) (*models.Weather, error) {

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", cityName, cfg.WeatherApiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get weather: %s", resp.Status)
	}

	var data Response

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &models.Weather{
		CityName:    cityName,
		Temperature: data.Main["temp"],
		FeelsLike:   data.Main["feels_like"],
		Pressure:    data.Main["pressure"],
		Humidity:    data.Main["humidity"],
	}, nil
}
