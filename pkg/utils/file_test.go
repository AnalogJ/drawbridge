package utils_test

import (
	"github.com/analogj/drawbridge/pkg/utils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

//func TestExpandPath(t *testing.T) {
//	t.Parallel()
//
//	//test
//	actual, err := utils.ExpandPath("~/test.file")
//
//	//assert
//	require.NoError(t, err,"should not raise an error")
//	require.Equal(t, []string{"example"}, actual, "should correctly retrieve keys from a map")
//}

func TestFileWrite(t *testing.T) {
	t.Parallel()

	//setup
	parentPath, err := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	testFilePath := filepath.Join(parentPath, "testfile.txt")

	//test
	err = utils.FileWrite(testFilePath, "test content", 0666, false)

	//assert
	require.NoError(t, err, "should not raise an error when writing file")
	require.FileExists(t, testFilePath, "test file should exist")
}

func TestFileWrite_DryRun(t *testing.T) {
	t.Parallel()

	//setup
	parentPath, err := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	testFilePath := filepath.Join(parentPath, "testfile_dryrun.txt")

	//test
	err = utils.FileWrite(testFilePath, "test content", 0666, true)

	//assert
	require.NoError(t, err, "should not raise an error when writing file")
	require.False(t, utils.FileExists(testFilePath), "test file should not be written in dry run mode")
}

func TestFileExists(t *testing.T) {
	t.Parallel()

	//test
	actual := utils.FileExists("testdata/placeholder.txt")

	//assert
	require.True(t, actual, "should detect that placeholder text file exists")
}

func TestFileExists_Invalid(t *testing.T) {
	t.Parallel()

	//test
	actual := utils.FileExists("testdata/placeholder-invalid.txt")

	//assert
	require.False(t, actual, "should detect that invalid file does not exist")
}

func TestFileDelete(t *testing.T) {
	t.Parallel()

	//setup
	parentPath, err := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	testFilePath := filepath.Join(parentPath, "testfile-delete.txt")

	//test
	err = utils.FileWrite(testFilePath, "test content", 0666, false)
	err = utils.FileDelete(testFilePath)

	//assert
	require.NoError(t, err, "should not raise an error when deleting file")
	require.False(t, utils.FileExists(testFilePath), "test file should not exist after deletion")
}
