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
		"active_custom_templates": []string{},
		"config": map[string]interface{}{
			"filepath": "/Users/jason/.ssh/drawbridge/test-app-idle-us-east-1",
			"pem_filepath": "/Users/jason/.ssh/drawbridge/pem/test/aws-test.pem",
		},
		"config_dir": "~/.ssh/drawbridge",
		"custom": []string{},
		"environment": "test",
		"pem_dir": "~/.ssh/drawbridge/pem",
		"shard": "us-east-1",
		"shard_type": "idle",
		"stack_name": "app",
		"ui_group_priority": []string{"environment","stack_name","shard","shard_type"},
		"ui_question_hidden": []string{},
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

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "valid_configfile_with_answers.yaml"))
	require.NoError(t, err, "should allow overriding default config template.")

	projList, err := project.CreateProjectListFromProvidedAnswers(testConfig)
	actualSortedList := projList.GetAll()

	//assert
	require.NoError(t, err, "should correctly get answers from config.")
	require.Equal(t, 5, projList.Length())
	require.Equal(t, 5, len(actualSortedList))
	require.Equal(t, actualSortedList[0], projList.GetIndex(0))
}
