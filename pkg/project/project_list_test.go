package project_test

import (
	"github.com/analogj/drawbridge/pkg/config"
	"github.com/analogj/drawbridge/pkg/project"
	"github.com/stretchr/testify/require"
	"path"
	"path/filepath"
	"testing"
)

func TestProjectList_WithEmptyAnswersList(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()
	err := testConfig.ReadConfig(path.Join("testdata", "valid_configfile_no_answers.yaml"))
	require.NoError(t, err, "should allow overriding default config template.")

	//test
	projList, err := project.CreateProjectListFromProvidedAnswers(testConfig)
	require.NoError(t, err, "should correctly load project list")

	actualSortedList := projList.GetAll()

	//assert
	require.Equal(t, 0, projList.Length(), "should correctly load provided answers")
	require.Equal(t, 0, len(actualSortedList), "should correcty populate sorted list after grouping")

	_, err = projList.GetIndex(0)
	require.Error(t, err, "should raise an error if attempting to access an empty list ProjectListEmptyError")

}

func TestProjectList_WithEmptyConfigDir(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()
	err := testConfig.ReadConfig(filepath.Join("testdata", "valid_configfile_no_answers.yaml"))
	require.NoError(t, err, "should allow overriding default config template.")
	testConfig.Set("options.config_dir", filepath.Join("testdata", "empty_config_dir"))

	//test
	projList, err := project.CreateProjectListFromConfigDir(testConfig)
	require.NoError(t, err, "should correctly load project list")

	actualSortedList := projList.GetAll()

	//assert
	require.Equal(t, 0, projList.Length(), "should correctly load provided answers")
	require.Equal(t, 0, len(actualSortedList), "should correcty populate sorted list after grouping")

	_, err = projList.GetIndex(0)
	require.Error(t, err, "should raise an error if attempting to access an empty list ProjectListEmptyError")

}

func TestProjectList_GetIndex(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()
	err := testConfig.ReadConfig(path.Join("testdata", "valid_configfile_with_answers.yaml"))
	require.NoError(t, err, "should allow overriding default config template.")

	//test
	projList, err := project.CreateProjectListFromProvidedAnswers(testConfig)
	require.NoError(t, err, "should correctly load project list")
	_, startErr := projList.GetIndex(0)
	_, lastErr := projList.GetIndex(4)
	_, lenthErr := projList.GetIndex(projList.Length())

	//assert
	require.Equal(t, 5, projList.Length(), "should correctly load provided answers")
	require.NoError(t, startErr, "should correctly retrieve item at start")
	require.NoError(t, lastErr, "should correctly retrieve item at end")
	require.Error(t, lenthErr, "should raise an error when accessing last item + 1 index")
}
