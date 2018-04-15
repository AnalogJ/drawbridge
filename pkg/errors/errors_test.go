package errors_test

import (
	"drawbridge/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

//func TestCheckErr_WithoutError(t *testing.T) {
//	t.Parallel()
//
//	//assert
//	require.NotPanics(t, func() {
//		errors.CheckErr(nil)
//	})
//}

//func TestCheckErr_Error(t *testing.T) {
//	t.Parallel()
//
//	//assert
//	require.Panics(t, func() {
//		errors.CheckErr(stderrors.New("This is an error"))
//	})
//}

func TestErrors(t *testing.T) {
	t.Parallel()

	//assert
	require.Implements(t, (*error)(nil), errors.EngineBuildPackageFailed("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.EngineBuildPackageInvalid("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.EngineDistCredentialsMissing("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.EngineDistPackageError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.EngineTestDependenciesError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.EngineTestRunnerError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.EngineTransformUnavailableStep("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.EngineUnspecifiedError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.EngineValidateToolError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.ScmAuthenticationFailed("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.ScmFilesystemError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.ScmPayloadFormatError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.ScmPayloadUnsupported("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.ScmUnauthorizedUser("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.ScmUnspecifiedError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.ScmCleanupFailed("test"), "should implement the error interface")
}
