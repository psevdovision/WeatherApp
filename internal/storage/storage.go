package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"weatherapp/internal/config"
)

type Weather struct {
	CreatedAt   string  `json:"-"`
	CityName    string  `json:"city_name"`
	Temperature float32 `json:"temp"`
	FeelsLike   float32 `json:"feels_like"`
	Pressure    float32 `json:"pressure"`
	Humidity    float32 `json:"humidity"`
}

type Storage struct {
	db *pgx.Conn
}

func New(cfg config.Config) (*Storage, error) {

	dsn := cfg.Dsn

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	storage := &Storage{db: conn}

	if err := storage.initializeDB(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %s", err)
	}

	return storage, nil
}

func (s *Storage) initializeDB() error {

	stmt := `
	CREATE TABLE IF NOT EXISTS weather (
		id SERIAL PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		city_name TEXT NOT NULL,
		temperature REAL NOT NULL,
		feels_like INT NOT NULL,
		pressure INT NOT NULL,
		humidity INT NOT NULL
	);`

	_, err := s.db.Exec(context.Background(), stmt)
	if err != nil {
		log.Printf("Error initializing database: %s", err)
		return err
	}

	log.Println("Database initialized successfully")

	return nil
}

func (s *Storage) SaveWeather(cityName string, temperature float32, feelsLike float32, pressure float32, humidity float32) error {

	stmt := `
	INSERT INTO weather (city_name, temperature, feels_like, pressure, humidity)
	VALUES ($1, $2, $3, $4, $5);`

	_, err := s.db.Exec(context.Background(), stmt, cityName, temperature, feelsLike, pressure, humidity)
	if err != nil {
		log.Printf("Failed to insert weather data: %s", err)
		return err
	}

	log.Printf("Weather data for city '%s' inserted successfully", cityName)
	return nil
}

func (s *Storage) GetWeather(cityName string) ([]Weather, error) {

	stmt := `SELECT city_name, temperature, feels_like, pressure, humidity
	FROM weather
	WHERE city_name = $1
	ORDER BY created_at DESC
	LIMIT 1;`

	rows, err := s.db.Query(context.Background(), stmt, cityName)
	if err != nil {
		log.Printf("Failed to get weather data: %s", err)
		return nil, err
	}
	defer rows.Close()

	var weatherList []Weather

	for rows.Next() {
		var weather Weather

		if err := rows.Scan(&weather.CityName, &weather.Temperature, &weather.FeelsLike, &weather.Pressure, &weather.Humidity); err != nil {
			log.Printf("Failed to scan row: %s", err)
			return nil, err
		}

		weatherList = append(weatherList, weather)
	}

	log.Printf("Weather data for city '%s' geted successfully", cityName)

	return weatherList, nil
}
