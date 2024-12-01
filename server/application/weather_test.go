package application_test

import (
	"github.com/rgoncalvesrr/fullcycle-labs-otel/server/application"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/server/pkg/weather"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWeather(t *testing.T) {
	celsius := weather.Celsius(0.0)
	expectedCelsius := float64(celsius)
	expectedFahrenheit := celsius.ToFahrenheit()
	expectedKelvin := celsius.ToKelvin()

	w, e := application.NewWeather(celsius)

	assert.Nil(t, e)
	assert.NotNil(t, w)
	assert.Equal(t, expectedCelsius, w.Celsius())
	assert.Equal(t, expectedKelvin, w.Kelvin())
	assert.Equal(t, expectedFahrenheit, w.Fahrenheit())
}

func TestNewWeatherShouldThrowError(t *testing.T) {
	w, e := application.NewWeather(-274)
	assert.Nil(t, w)
	assert.NotNil(t, e)
	assert.Equal(t, "temperature cannot be less than 273.15", e.Error())
	assert.Equal(t, application.ErrInvalidTemperature, e)
}
