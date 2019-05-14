package utils_test

import (
	"github.com/analogj/drawbridge/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

var snakeCaseTests = []struct {
	n        string // input
	expected string // expected result
}{
	{"this_is_an_input", "ThisIsAnInput"},
	{"", ""},
	{"hello", "Hello"},
}

func TestSnakeCaseToCamelCase(t *testing.T) {
	t.Parallel()
	for _, tt := range snakeCaseTests {
		//test
		actual := utils.SnakeCaseToCamelCase(tt.n)

		//assert
		require.Equal(t, tt.expected, actual, "should convert to camel case correctly")
	}
}

func TestStringToInt(t *testing.T) {
	t.Parallel()

	//test
	convInt, err := utils.StringToInt("1")

	//assert
	require.NoError(t, err)
	require.Equal(t, 1, convInt, "should correctly parse string")
}

func TestStringToInt_InvalidString(t *testing.T) {
	t.Parallel()

	//test
	_, err := utils.StringToInt("abc")

	//assert
	require.Error(t, err)
}

func TestStringToInt_EmptyString(t *testing.T) {
	t.Parallel()

	//test
	_, err := utils.StringToInt("")

	//assert
	require.Error(t, err)
}

func TestStripIndent(t *testing.T) {
	t.Parallel()

	testString := `
	this is my multi line string
	line2
	line 3`

	require.Equal(t, "\nthis is my multi line string\nline2\nline 3", utils.StripIndent(testString))
}
func TestLeftPad2Len(t *testing.T) {
	t.Parallel()

	require.Equal(t, "-----12345", utils.LeftPad2Len("12345", "-", 10))
	require.Equal(t, "345", utils.LeftPad2Len("12345", "-", 3))
}

func TestRightPad2Len(t *testing.T) {
	t.Parallel()

	require.Equal(t, "12345-----", utils.RightPad2Len("12345", "-", 10))
	require.Equal(t, "123", utils.RightPad2Len("12345", "-", 3))
}

func TestLeftPad(t *testing.T) {
	t.Parallel()

	require.Equal(t, "----------12345", utils.LeftPad("12345", "-", 10))
	require.Equal(t, "---12345", utils.LeftPad("12345", "-", 3))
}

func TestRightPad(t *testing.T) {
	t.Parallel()

	require.Equal(t, "12345----------", utils.RightPad("12345", "-", 10))
	require.Equal(t, "12345---", utils.RightPad("12345", "-", 3))
}
