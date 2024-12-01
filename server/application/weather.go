package application

import (
	"github.com/rgoncalvesrr/fullcycle-labs-otel/server/pkg/weather"
)

type Weather struct {
	celsius    weather.Celsius
	fahrenheit float64
	kelvin     float64
}

func NewWeather(tempCelsius weather.Celsius) (*Weather, error) {
	w := &Weather{
		celsius:    tempCelsius,
		kelvin:     tempCelsius.ToKelvin(),
		fahrenheit: tempCelsius.ToFahrenheit(),
	}

	if e := w.validate(); e != nil {
		return nil, e
	}

	return w, nil
}

func (w *Weather) Celsius() float64 {
	return float64(w.celsius)
}

func (w *Weather) Fahrenheit() float64 {
	return w.fahrenheit
}

func (w *Weather) Kelvin() float64 {
	return w.kelvin
}

func (w *Weather) validate() error {
	if w.celsius < -273.15 {
		return ErrInvalidTemperature
	}

	return nil
}
