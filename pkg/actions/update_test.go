package actions_test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"drawbridge/pkg/actions"
)

func TestUpdateAction_GetLatestReleaseInfo(t *testing.T) {
	t.Parallel()

	//setup
	updateAction := actions.UpdateAction{}

	//test
	releaseInfo, err := updateAction.GetLatestReleaseInfo()


	//assert
	require.NoError(t, err, "should not raise an error when retrieving release info")
	require.Equal(t, 3 , len(releaseInfo.Assets), "should correctly retrieve download info for 3 binaries")
}