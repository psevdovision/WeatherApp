package models

type Weather struct {
	CityName    string  `json:"city_name"`
	Temperature float32 `json:"temp"`
	Pressure    float32 `json:"pressure"`
	Humidity    float32 `json:"humidity"`
	FeelsLike   float32 `json:"feels_like"`
}
