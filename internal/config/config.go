package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Dsn           string     `yaml:"dsn" env:"DSN"`
	WeatherApiKey string     `yaml:"weather_api_key" env:"WEATHER_API_KEY"`
	HTTPServer    HTTPServer `yaml:"http_server" env-default:":8080"`
}

type HTTPServer struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout" end-default:"10s"`
}

func MustLoad() *Config {

	var cfg Config

	configPath := "your config path"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}

	return &cfg

}
