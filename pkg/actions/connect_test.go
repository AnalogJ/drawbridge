package actions_test

import (
	"github.com/analogj/drawbridge/pkg/actions"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

//func TestConnectAction_SshAgentAddPemKey(t *testing.T) {
//	t.Parallel()
//
//	//setup
//	connectAction := actions.ConnectAction{}
//
//	//test
//	err := connectAction.SshAgentAddPemKey(path.Join("testdata", "connect/test_rsa.pem"))
//
//	//assert
//	require.NoError(t, err, "should not raise an error when adding pem key to ssh-agent")
//}

func TestConnectAction_SshAgentAddPemKey_InvalidPath(t *testing.T) {
	t.Parallel()

	//setup
	connectAction := actions.ConnectAction{}

	//test
	err := connectAction.SshAgentAddPemKey(filepath.Join("testdata", "invalid_path.pem"))

	//assert
	require.Error(t, err, "should raise an error when adding invalid pem key to ssh-agent")
}
