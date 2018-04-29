package actions_test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"drawbridge/pkg/actions"
	"drawbridge/pkg/config"
	"path"
	"drawbridge/pkg/utils"
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
	err = utils.CopyDir(path.Join("testdata", "delete"), parentPath)
	require.NoError(t, err, "should not raise an error when deleting answer file")

	configData.Set("options.config_dir", parentPath)
	configData.Set("config_templates.default.pem_filepath", "test_rsa.pem")
	configData.Set("options.pem_dir", parentPath)
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
	}, true)


	//assert
	require.NoError(t, err, "should not raise an error when deleting answer file")
	require.False(t, utils.FileExists(path.join(parentPath, "prod-app-idle-us-east-1")), "test file should not be exist")

}
