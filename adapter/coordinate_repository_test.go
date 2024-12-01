package adapter_test

import (
	"context"
	"github.com/rgoncalvesrr/fullcycle-labs-otel/adapter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCoordinateRepository(t *testing.T) {
	r := adapter.NewCoordinateRepository()
	c, e := r.GetByCep(context.Background(), "09130220")

	assert.NotNil(t, r)
	assert.Nil(t, e)
	assert.NotNil(t, c)
	assert.Equal(t, "Santo André", c.City)
	assert.Equal(t, "-23.6929", c.Latitude)
	assert.Equal(t, "-46.50423", c.Longitude)

}
