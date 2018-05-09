package actions_test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"drawbridge/pkg/actions"
	"drawbridge/pkg/config"
	"io/ioutil"
	"os"
)

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
