package updater

import (
	"log"
	"time"
	"weatherapp/internal/api_client"
	"weatherapp/internal/config"
	"weatherapp/internal/storage"
)

func WeatherUpdater(db *storage.Storage, cfg config.Config) {

	cities := []string{"Moscow", "Irkutsk", "Ufa", "Tokyo", "London"}

	go func() {
		for {
			for _, city := range cities {
				weather, err := api_client.GetWeatherByURL(city, cfg)
				if err != nil {
					log.Printf("Failed to get weather for %s:%s", city, err)
					continue
				}

				err = db.SaveWeather(weather.CityName, weather.Temperature, weather.FeelsLike, weather.Pressure, weather.Humidity)
				if err != nil {
					log.Printf("Failed to save weather for %s:%s", city, err)
				}
			}

			time.Sleep(10 * time.Minute)
		}
	}()
}
