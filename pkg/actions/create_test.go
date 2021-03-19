package actions_test

import (
	"github.com/analogj/drawbridge/pkg/actions"
	"github.com/analogj/drawbridge/pkg/config"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateAction_Start_AnswersFileWithOverrides(t *testing.T) {
	t.Parallel()

	//setup
	configData, err := config.Create()
	require.NoError(t, err)
	//read test config file w/Answers
	err = configData.ReadConfig(filepath.Join("testdata", "create", "valid_answers_override_active_custom_template.yaml"))

	parentPath, err := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)

	configData.Set("options.config_dir", parentPath)
	configData.Set("options.pem_dir", parentPath)
	configData.Set("options.active_config_template", "override")
	configData.Set("options.active_custom_templates", []string{})
	createAction := actions.CreateAction{
		Config: configData,
	}

	//test
	err = createAction.Start(map[string]interface{}{
		"environment": "test",
		"stack_name":  "tested",
		"shard":       "us-east-1",
		"shard_type":  "live",
		"username":    "aws",
	}, false)

	//assert
	require.NoError(t, err, "should not raise an error when adding writing answer file")
	require.FileExists(t, filepath.Join(parentPath, "aws"))
	require.NoFileExists(t, filepath.Join(parentPath, "custom-template-test-aws"))
}

func TestCreateAction_WriteAnswersFile(t *testing.T) {
	t.Parallel()

	//setup
	configData, err := config.Create()
	require.NoError(t, err)

	parentPath, err := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)

	configData.Set("options.config_dir", parentPath)
	createAction := actions.CreateAction{
		Config: configData,
	}

	//test
	err = createAction.WriteAnswersFile("base_name", map[string]interface{}{
		"test": "test_data",
	}, false)

	//assert
	require.NoError(t, err, "should not raise an error when adding writing answer file")
}
