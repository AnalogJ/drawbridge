package actions_test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"drawbridge/pkg/actions"
	"drawbridge/pkg/config"
	"drawbridge/pkg/utils"
	"path/filepath"
	"path"
)

func patchEnv(key, value string) func() {
	bck := os.Getenv(key)
	deferFunc := func() {
		os.Setenv(key, bck)
	}

	os.Setenv(key, value)
	return deferFunc
}

func TestDeleteAction_One(t *testing.T) {
	t.Parallel()

	//setup
	configData, err := config.Create()
	require.NoError(t, err)

	parentPath, err := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	defer patchEnv("HOME", parentPath)()
	drawbridgePath := path.Join(parentPath, "drawbridge")
	err = utils.CopyDir(path.Join("testdata", "delete"), drawbridgePath)
	require.NoError(t, err, "should not raise an error when deleting answer file")

	configData.Set("options.config_dir", drawbridgePath)
	configData.Set("config_templates.default.pem_filepath", "test_rsa.pem")
	configData.Set("options.pem_dir", drawbridgePath)
	deleteAction := actions.DeleteAction{
		Config: configData,
	}

	//test
	err = deleteAction.One(map[string]interface{}{
		"environment": "prod",
		"stack_name": "app",
		"shard": "us-east-1",
		"shard_type": "idle",
		"username": "aws",
		"config": map[string]interface{}{
			"filepath": path.Join(drawbridgePath, "prod-app-idle-us-east-1"),
		},
		"config_dir": drawbridgePath,
	}, true)


	//assert
	require.NoError(t, err, "should not raise an error when deleting answer file")
	require.False(t, utils.FileExists(filepath.Join(drawbridgePath, "prod-app-idle-us-east-1")), "test file should not be exist")

}

func TestDeleteAction_All(t *testing.T) {
	t.Parallel()

	//setup
	configData, err := config.Create()
	require.NoError(t, err)

	parentPath, err := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	defer patchEnv("HOME", parentPath)()
	drawbridgePath := path.Join(parentPath, "drawbridge")
	err = utils.CopyDir(path.Join("testdata", "delete"), drawbridgePath)
	require.NoError(t, err, "should not raise an error when deleting answer file")

	configData.Set("options.config_dir", drawbridgePath)
	configData.Set("config_templates.default.pem_filepath", "test_rsa.pem")
	configData.Set("options.pem_dir", drawbridgePath)
	deleteAction := actions.DeleteAction{
		Config: configData,
	}

	//test
	err = deleteAction.All([]map[string]interface{}{
		"environment": "prod",
		"stack_name": "app",
		"shard": "us-east-1",
		"shard_type": "idle",
		"username": "aws",
		"config": map[string]interface{}{
			"filepath": path.Join(drawbridgePath, "prod-app-idle-us-east-1"),
		},
		"config_dir": drawbridgePath,
	}, true)


	//assert
	require.NoError(t, err, "should not raise an error when deleting answer file")
	require.False(t, utils.FileExists(filepath.Join(drawbridgePath, "prod-app-idle-us-east-1")), "test file should not be exist")

}
