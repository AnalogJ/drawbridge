package utils_test

import (
	"github.com/analogj/drawbridge/pkg/utils"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestBashCmdExec(t *testing.T) {
	t.Parallel()

	//test
	cerr := utils.BashCmdExec("echo 'hello from bash'", "", nil, "")

	//assert
	require.NoError(t, cerr)
}

func TestBashCmdExec_StdErr(t *testing.T) {
	t.Parallel()

	//test
	cerr := utils.BashCmdExec("(>&2 echo 'test writing to stderr')", "", nil, "")

	//assert
	require.NoError(t, cerr)
}

func TestBashCmdExec_Prefix(t *testing.T) {
	t.Parallel()

	//test
	cerr := utils.BashCmdExec("echo 'hello from bash with custom prefix'", "", nil, "cust_prefix")

	//assert
	require.NoError(t, cerr)
}

func TestCmdExec_Date(t *testing.T) {
	t.Parallel()

	//test
	cerr := utils.CmdExec("date", []string{}, "", nil, "")

	//assert
	require.NoError(t, cerr)
}

func TestCmdExec_Echo(t *testing.T) {
	t.Parallel()

	//test
	cerr := utils.CmdExec("echo", []string{"hello", "world"}, "", nil, "")

	//assert
	require.NoError(t, cerr)
}

func TestCmdExec_Error(t *testing.T) {
	t.Skip()
	//t.Parallel()

	//test
	cerr := utils.CmdExec("sh", []string{"-c", "'exit 1'"}, "", nil, "")

	//assert
	require.Error(t, cerr)
}

func TestCmdExec_WorkingDirRelative(t *testing.T) {
	t.Parallel()

	//test
	cerr := utils.CmdExec("ls", []string{}, "testdata", nil, "")

	//assert
	require.Error(t, cerr)
}

func TestCmdExec_WorkingDirAbsolute(t *testing.T) {
	t.Parallel()

	//test
	absPath, aerr := filepath.Abs(".")
	cerr := utils.CmdExec("ls", []string{}, absPath, nil, "")

	//assert
	require.NoError(t, aerr)
	require.NoError(t, cerr)
}
