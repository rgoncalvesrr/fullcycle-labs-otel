package configs_test

import (
	"github.com/rgoncalvesrr/fullcycle-labs-otel/configs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	c := configs.Cfg
	assert.NotNil(t, c)
}
