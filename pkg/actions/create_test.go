package actions_test

import (
	"github.com/analogj/drawbridge/pkg/actions"
	"github.com/analogj/drawbridge/pkg/config"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
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
