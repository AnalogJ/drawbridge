package utils_test

import (
	"drawbridge/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPopulateTemplate(t *testing.T) {
	t.Parallel()

	//test
	actual, err := utils.PopulateTemplate("test {{.example}}", map[string]interface{}{"example": "17"})

	//assert
	require.NoError(t, err, "should not throw an error")
	require.Equal(t, "test 17", actual, "should populate a template correctly")
}

func TestPopulateTemplate_MissingDataShouldReturnErr(t *testing.T) {
	t.Parallel()

	//test
	_, err := utils.PopulateTemplate("test {{.example}}", map[string]interface{}{"example1": "17"})

	//assert
	require.Error(t, err, "should throw an error if missing template data")
}

func TestPopulateTemplate_SliceData(t *testing.T) {
	t.Parallel()

	//test
	actual, err := utils.PopulateTemplate("{{range .}}test {{.example1}},{{end}}", []map[string]interface{}{
		{"example1": "17"},
		{"example1": "18"},
	})

	//assert
	require.NoError(t, err, "should not throw an error")
	require.Equal(t, "test 17,test 18,", actual, "should populate a template correctly")
}

func TestPopulateTemplate_UniquePort(t *testing.T) {
	t.Parallel()

	//test
	str, err := utils.PopulateTemplate("test {{uniquePort .}}", map[string]interface{}{"example1": "17"})

	//assert
	require.NoError(t, err, "should throw an error if missing template data")
	require.Equal(t, "test 48275", str, "should correctly popualte unique port")
}

func TestUniquePort(t *testing.T) {
	t.Parallel()

	//test
	port, err := utils.UniquePort(map[string]interface{}{"example1": "17"})

	//assert
	require.NoError(t, err, "should not raise an error")
	require.Equal(t, 48275, port, "should generate repeatible unique port from data")
}

func TestUniquePort_WithNonStringValues(t *testing.T) {
	t.Parallel()

	//test
	port, err := utils.UniquePort(map[string]interface{}{"example1": "17", "example2": 18})

	//assert
	require.NoError(t, err, "should not raise an error")
	require.Equal(t, 38792, port, "should generate repeatible unique port from data")
}

func TestUniquePort_WithStringKey(t *testing.T) {
	t.Parallel()

	//test
	port, err := utils.UniquePort("17")

	//assert
	require.NoError(t, err, "should not raise an error")
	require.Equal(t, 36016, port, "should generate repeatible unique port from data")
}
