package adapter

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/server/application"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/server/configs"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/server/pkg/weather"
)

type WeatherApiOutput struct {
	Current struct {
		Celsius float64 `json:"temp_c"`
	}
}

type weatherRepository struct {
}

func NewWeatherRepository() application.IWeatherRepository {

	return &weatherRepository{}
}

func (w *weatherRepository) GetTemperature(ctx context.Context, coordinate *application.Coordinate) (*application.Weather, error) {
	client := resty.New()
	r, e := client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetQueryParams(map[string]string{
			"key": configs.Cfg.WeatherApiKey,
			"q":   fmt.Sprintf("%s,%s", coordinate.Latitude, coordinate.Longitude),
		}).
		SetResult(&WeatherApiOutput{}).
		Get(configs.Cfg.WeatherApiUrl)
	if e != nil {
		return nil, e
	}

	result, e := application.NewWeather(weather.Celsius(r.Result().(*WeatherApiOutput).Current.Celsius))

	if e != nil {
		return nil, e
	}

	return result, nil
}
