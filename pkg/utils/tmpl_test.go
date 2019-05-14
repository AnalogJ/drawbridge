package utils_test

import (
	"github.com/analogj/drawbridge/pkg/utils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func patchEnv(key, value string) func() {
	bck := os.Getenv(key)
	deferFunc := func() {
		os.Setenv(key, bck)
	}

	os.Setenv(key, value)
	return deferFunc
}

func TestPopulatePathTemplate(t *testing.T) {
	t.Parallel()

	//test
	actual, err := utils.PopulatePathTemplate("/tmp/{{.example}}", map[string]interface{}{"example": "17"})

	//assert
	require.NoError(t, err, "should not throw an error")
	require.Equal(t, "/tmp/17", actual, "should populate a template correctly")
}

func TestPopulatePathTemplate_JoinedPath(t *testing.T) {
	t.Parallel()

	//test
	actual, err := utils.PopulatePathTemplate(path.Join("/tmp", "{{.example}}"), map[string]interface{}{"example": "17"})

	//assert
	require.NoError(t, err, "should not throw an error")
	require.Equal(t, "/tmp/17", actual, "should populate a template correctly")
}

func TestPopulatePathTemplate_RelativePath(t *testing.T) {
	t.Parallel()

	//setup
	parentPath, err := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	defer patchEnv("HOME", parentPath)()

	//test
	actual, err := utils.PopulatePathTemplate(path.Join("~/", "{{.example}}"), map[string]interface{}{"example": "17"})

	//assert
	require.NoError(t, err, "should not throw an error")
	require.Equal(t, path.Join(parentPath, "17"), actual, "should populate a template correctly")
}

func TestPopulateTemplate(t *testing.T) {
	t.Parallel()

	//test
	actual, err := utils.PopulateTemplate("test {{.example}}", map[string]interface{}{"example": "17"})

	//assert
	require.NoError(t, err, "should not throw an error")
	require.Equal(t, "test 17", actual, "should populate a template correctly")
}

func TestPopulateTemplate_InvalidTemplate(t *testing.T) {
	t.Parallel()

	//test
	_, err := utils.PopulateTemplate("test {{.example", map[string]interface{}{"example": "17"})

	//assert
	require.Error(t, err, "should throw an error")
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

func TestPopulateTemplate_StringsHasPrefix(t *testing.T) {
	t.Parallel()

	//test
	str, err := utils.PopulateTemplate(`test {{if stringsHasPrefix .example1 "this-is"}}{{stringsTrimPrefix .example1 "this-is-a-test-"}}{{end}}`, map[string]interface{}{"example1": "this-is-a-test-string"})

	//assert
	require.NoError(t, err, "should throw an error if missing template data")
	require.Equal(t, "test string", str, "should correctly test for prefix, and trim prefix")
}
