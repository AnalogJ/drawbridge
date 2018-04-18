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
	require.Implements(t, (*error)(nil), errors.ConfigFileMissingError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.ConfigValidationError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.TemplateFileExistsError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.QuestionKeyInvalidError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.QuestionValidationError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.AnswerFormatError("test"), "should implement the error interface")
}
