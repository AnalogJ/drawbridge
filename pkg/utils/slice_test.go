package utils_test

import (
	"testing"
	"drawbridge/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestSliceIncludes(t *testing.T) {
	t.Parallel()

	//test
	actual := utils.SliceIncludes([]string{"example", "example2", "example3"}, "example")

	//assert
	require.True(t, actual,"should find item in slice")
}

func TestSliceIncludes_WithInvalid(t *testing.T) {
	t.Parallel()

	//test
	actual := utils.SliceIncludes([]string{"example", "example2", "example3"}, "nothere")

	//assert
	require.False(t, actual,"should not find item in slice")
}