package template_test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"drawbridge/pkg/config/template"
	"os"
	"path"
	"drawbridge/pkg/utils"
	"io/ioutil"
)

func TestFileTemplate_PopulateFilePath(t *testing.T) {
	t.Parallel()

	//setup
	fileTemplate := template.FileTemplate{
		FilePath: "/{{.example}}/{{.example2}}.text",
		Template: template.Template{
			Content: "",
		},
	}

	//test
	actual, err := fileTemplate.PopulateFilePath(map[string]interface{}{
		"example": "1",
		"example2": "2",
	})

	//assert
	require.NoError(t, err,"should not raise an error when populating filepath template")
	require.Equal(t, "/1/2.text", actual, "should correctly populate template")
}

func TestFileTemplate_PopulateFilePath_WithMissingData(t *testing.T) {
	t.Parallel()

	//setup
	fileTemplate := template.FileTemplate{
		FilePath: "/{{.example}}/{{.example2}}.text",
		Template: template.Template{
			Content: "",
		},
	}

	//test
	_, err := fileTemplate.PopulateFilePath(map[string]interface{}{
		"example": "1",
	})

	//assert
	require.Error(t, err,"should raise an error when populating filepath template")
}

func TestFileTemplate_DeleteTemplate(t *testing.T) {
	t.Parallel()

	//setup
	parentPath, err := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)

	testFilePathTemplate := path.Join(parentPath, "{{.example}}.text")
	testFilePath := path.Join(parentPath, "1.text")
	err = utils.FileWrite(testFilePath, "test content", 0666, false)
	require.NoError(t, err, "should not raise an error when writing test file.")

	fileTemplate := template.FileTemplate{
		FilePath: testFilePathTemplate,
		Template: template.Template{
			Content: "",
		},
	}

	//test
	err = fileTemplate.DeleteTemplate(map[string]interface{}{
		"example": "1",
	})

	//assert
	require.NoError(t, err,"should not raise an error deleting filepath template")
	require.False(t, utils.FileExists(testFilePath), "test file should not be exist")
}

func TestFileTemplate_DeleteTemplate_WhenFileDoesNotExist(t *testing.T) {
	t.Parallel()

	//setup
	parentPath, err := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)

	testFilePathTemplate := path.Join(parentPath, "{{.example}}/{{.example2}}.text")

	fileTemplate := template.FileTemplate{
		FilePath: testFilePathTemplate,
		Template: template.Template{
			Content: "",
		},
	}

	//test
	err = fileTemplate.DeleteTemplate(map[string]interface{}{
		"example": "1",
		"example2": "2",
	})

	//assert
	require.NoError(t, err,"should not raise an error deleting filepath template")
}

func TestFileTemplate_WriteTemplate(t *testing.T) {
	t.Parallel()

	//setup
	parentPath, err := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)

	testFilePathTemplate := path.Join(parentPath, "{{.example}}/{{.example2}}.text")
	testFilePath := path.Join(parentPath, "1/2.text")

	fileTemplate := template.FileTemplate{
		FilePath: testFilePathTemplate,
		Template: template.Template{
			Content: "{{.content}}",
		},
	}

	//test
	actual, err := fileTemplate.WriteTemplate(map[string]interface{}{
		"example": "1",
		"example2": "2",
		"content": "this is my content",
	}, false)

	//assert
	require.NoError(t, err,"should not raise an error deleting filepath template")
	require.FileExists(t, testFilePath, "should write file to correct path")
	require.Equal(t, map[string]interface{}{"filepath": testFilePath}, actual, "should return some metadata about the template")
}

func TestFileTemplate_WriteTemplate_WhenDestinationExists(t *testing.T) {
	t.Parallel()

	//setup
	parentPath, err := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)

	testFilePathTemplate := path.Join(parentPath, "{{.example}}.text")
	testFilePath := path.Join(parentPath, "1.text")
	err = utils.FileWrite(testFilePath, "previous content", 0666, false)
	require.NoError(t, err, "should not raise an error when writing test file.")

	fileTemplate := template.FileTemplate{
		FilePath: testFilePathTemplate,
		Template: template.Template{
			Content: "{{.content}}",
		},
	}

	//test
	_, err = fileTemplate.WriteTemplate(map[string]interface{}{
		"example": "1",
		"content": "this is my content",
	}, false)

	//assert
	require.Error(t, err,"should raise an error if destination file already exists.")
}