package actions_test

import (
	"github.com/analogj/drawbridge/pkg/actions"
	"github.com/analogj/drawbridge/pkg/config"
	"github.com/analogj/drawbridge/pkg/utils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestProxyAction_Start(t *testing.T) {
	t.Parallel()

	//setup
	configData, err := config.Create()
	require.NoError(t, err)

	parentPath, err := ioutil.TempDir("", "")
	defer os.RemoveAll(parentPath)
	defer patchEnv("HOME", parentPath)()
	drawbridgePath := filepath.Join(parentPath, "drawbridge")
	err = utils.CopyDir(filepath.Join("testdata", "delete"), drawbridgePath)
	require.NoError(t, err, "should not raise an error when deleting answer file")

	configData.Set("options.config_dir", drawbridgePath)
	configData.Set("config_templates.default.pem_filepath", "test_rsa.pem")
	configData.Set("options.pem_dir", drawbridgePath)
	configData.Set("pac_template.filepath", filepath.Join(parentPath, "drawbridge.pac"))

	proxyAction := actions.ProxyAction{
		Config: configData,
	}

	//test
	err = proxyAction.Start([]map[string]interface{}{
		{
			"environment": "prod",
			"stack_name":  "app",
			"shard":       "us-east-1",
			"shard_type":  "idle",
			"username":    "aws",
			"config": map[string]interface{}{
				"filepath": filepath.Join(drawbridgePath, "prod-app-idle-us-east-1"),
			},
			"config_dir": drawbridgePath,
		},
	}, false)

	//assert
	require.NoError(t, err, "should not raise an error when generating pac file")
	require.FileExists(t, filepath.Join(parentPath, "drawbridge.pac"))
}
