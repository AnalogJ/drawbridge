package utils_test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"path"
	"os"
	"bytes"
	"gopkg.in/yaml.v2"
	"drawbridge/pkg/utils"
)

func TestStringifyYAMLMapKeys(t *testing.T) {
	t.Parallel()

	//setup
	testData, err := os.Open(path.Join("testdata", "test.yaml"))
	require.NoError(t, err, "should not throw an error")

	buf := new(bytes.Buffer)
	buf.ReadFrom(testData)
	parsedMap := map[interface{}]interface{}{}
	err = yaml.Unmarshal(buf.Bytes(), &parsedMap)
	require.NoError(t, err, "should not throw an error")

	//test
	stringifiedMap := utils.StringifyYAMLMapKeys(parsedMap)

	//assert
	require.Equal(t, map[string]interface{}{
		"test_key": "value",
		"test_number": 1,
		"test_nested": map[string]interface{}{
			"test_level_1": "hellp",
		},
	}, stringifiedMap, "should correctly stringify map")
}
