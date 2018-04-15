package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfiguration(t *testing.T) {

	//test
	config := new(configuration)

	//assert
	require.Implements(t, (*Interface)(nil), config, "should implement the config interface")
}
