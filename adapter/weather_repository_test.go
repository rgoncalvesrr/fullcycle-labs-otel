package adapter_test

import (
	"context"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/adapter"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/application"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWeatherRepository(t *testing.T) {
	cord := application.NewCoordinate("Santo André", "-23.667", "-46.517")
	r := adapter.NewWeatherRepository()
	c, e := r.GetTemperature(context.Background(), cord)
	assert.Nil(t, e)
	assert.NotNil(t, r)
	assert.NotNil(t, c)
}
