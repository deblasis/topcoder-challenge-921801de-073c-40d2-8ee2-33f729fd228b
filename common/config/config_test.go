package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	config, e := LoadConfig()
	require.Nil(t, e)
	require.NotNil(t, config)
}
