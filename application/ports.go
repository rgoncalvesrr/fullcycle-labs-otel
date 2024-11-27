package application

import "context"

type ICoordinateRepository interface {
	GetByCep(ctx context.Context, cep string) (*Coordinate, error)
}

type IWeatherRepository interface {
	GetTemperature(ctx context.Context, coordinate *Coordinate) (*Weather, error)
}
