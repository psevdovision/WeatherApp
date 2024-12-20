package main

import (
	"log"
	"net/http"
	"weatherapp/internal/config"
	"weatherapp/internal/http-server/router"
	"weatherapp/internal/storage"
	"weatherapp/internal/updater"
)

func main() {

	log.Println("Starting WeatherApp")

	cfg := config.MustLoad()

	db, err := storage.New(*cfg)
	if err != nil {
		log.Fatal("Failed to connect to database in main")
	}

	updater.WeatherUpdater(db, *cfg)

	server := router.New(cfg, db)

	port := cfg.HTTPServer.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server started on port %s", port)

	if err := http.ListenAndServe("0.0.0.0:3030", server); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
