package adapter

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/configs"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/internal/application"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/pkg/weather"
	"go.opentelemetry.io/otel/trace"
)

type WeatherApiOutput struct {
	Current struct {
		Celsius float64 `json:"temp_c"`
	}
}

type weatherRepository struct {
	tracer trace.Tracer
}

func NewWeatherRepository(tracer trace.Tracer) application.IWeatherRepository {
	return &weatherRepository{
		tracer: tracer,
	}
}

func (w *weatherRepository) GetTemperature(ctx context.Context, coordinate *application.Coordinate) (*application.Weather, error) {
	ctx, span := w.tracer.Start(ctx, "get-external-weather")
	defer span.End()

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
