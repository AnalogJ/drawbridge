package project_test

import (
	"testing"
	"path"
	"github.com/stretchr/testify/require"
	"drawbridge/pkg/project"
	"drawbridge/pkg/config"
)

func TestCreateProjectFromConfigDirAnswerFile(t *testing.T) {
	t.Parallel()

	//test
	answerFile := path.Join("testdata", "valid_answerfile.yaml")
	proj, err := project.CreateProjectFromConfigDirAnswerFile(answerFile)
	require.NoError(t, err, "should correctly parse answerfile.")

	//assert
	require.Equal(t, map[string]interface{}{
		"active_config_template": "default",
		"active_custom_templates": []interface{}{},
		"config": map[string]interface{}{
			"filepath": "/Users/jason/.ssh/drawbridge/test-app-idle-us-east-1",
			"pem_filepath": "/Users/jason/.ssh/drawbridge/pem/test/aws-test.pem",
		},
		"config_dir": "~/.ssh/drawbridge",
		"custom": []interface{}{},
		"environment": "test",
		"pem_dir": "~/.ssh/drawbridge/pem",
		"shard": "us-east-1",
		"shard_type": "idle",
		"stack_name": "app",
		"ui_group_priority": []interface{}{"environment","stack_name","shard","shard_type"},
		"ui_question_hidden": []interface{}{},
		"username":"aws",
	}, proj.Answers, "should parse populate")
	require.Equal(t, answerFile, proj.AnswerFilePath, "correctly set the answerfile path")
	require.Equal(t, "/Users/jason/.ssh/drawbridge/test-app-idle-us-east-1", proj.ConfigFilePath, "correctly set the config filepath")
	require.Equal(t, "/Users/jason/.ssh/drawbridge/pem/test/aws-test.pem", proj.PemFilePath, "correctly set the pem filepath")
}

func TestCreateProjectListFromProvidedAnswers(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()
	err := testConfig.ReadConfig(path.Join("testdata", "valid_configfile_with_answers.yaml"))
	require.NoError(t, err, "should allow overriding default config template.")

	//test
	projList, err := project.CreateProjectListFromProvidedAnswers(testConfig)
	require.NoError(t, err, "should correctly load project list")

	actualSortedList := projList.GetAll()
	actualFirstAnswer, err := projList.GetIndex(0)
	require.NoError(t, err, "should correctly get item at index")


	//assert
	require.NoError(t, err, "should correctly get answers from config.")
	require.Equal(t, 5, projList.Length(), "should correctly load provided answers")
	require.Equal(t, 5, len(actualSortedList), "should correcty populate sorted list after grouping")
	require.Equal(t, actualSortedList[0], actualFirstAnswer, "total list lenth provided should match list length after grouping")
	require.Equal(t, map[string]interface {}{
		"environment":"test",
		"stack_name":"test2",
		"shard":"us-east-1",
		"shard_type":"live",
		"username":"aws",
	}, actualFirstAnswer, "sort order should always be consistent")

}


func TestCreateProjectListFromConfigDir(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()
	err := testConfig.ReadConfig(path.Join("testdata", "valid_configfile_with_answers.yaml"))
	require.NoError(t, err, "should allow overriding default config template.")
	testConfig.Set("options.config_dir", path.Join("testdata", "config_dir"))


	//test
	projList, err := project.CreateProjectListFromConfigDir(testConfig)
	require.NoError(t, err, "should correctly load project list")

	actualSortedList := projList.GetAll()
	actualFirstAnswer, err := projList.GetIndex(0)
	require.NoError(t, err, "should correctly get item at index")


	//assert
	require.NoError(t, err, "should correctly get answers from config directory")
	require.Equal(t, 9, projList.Length(), "should correctly load provided answers")
	require.Equal(t, 9, len(actualSortedList), "should correcty populate sorted list after grouping")
	require.Equal(t, actualSortedList[0], actualFirstAnswer, "total list lenth provided should match list length after grouping")
	require.Equal(t, map[string]interface {}{
		"ui_question_hidden":[]interface {}{},
		"active_config_template":"default",
		"active_custom_templates":[]interface {}{},
		"config":map[string]interface {}{
			"filepath":"/Users/jason/.ssh/drawbridge/prod-app-idle-us-east-1",
			"pem_filepath":"/Users/jason/.ssh/drawbridge/pem/prod/aws-prod.pem",
		},
		"environment":"prod",
		"pem_dir":"~/.ssh/drawbridge/pem",
		"shard_type":"idle",
		"stack_name":"app",
		"username":"aws",
		"config_dir":"~/.ssh/drawbridge",
		"custom":[]interface {}{},
		"shard":"us-east-1",
		"ui_group_priority":[]interface {}{"environment", "stack_name", "shard", "shard_type"},
	}, actualFirstAnswer, "sort order should always be consistent")
}