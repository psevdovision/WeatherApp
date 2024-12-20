package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"weatherapp/internal/config"
	"weatherapp/internal/http-server/handlers"
	"weatherapp/internal/storage"
)

func New(cfg *config.Config, storage *storage.Storage) http.Handler {

	r := mux.NewRouter()

	weatherHandler := handlers.NewWeatherHandler(storage, *cfg)

	r.HandleFunc("/api/weather/getByCityName", weatherHandler.GetByCityName).Methods(http.MethodPost)

	r.HandleFunc("/api/weather/getPopularCities", weatherHandler.GetPopularCities).Methods(http.MethodGet)

	return r
}
