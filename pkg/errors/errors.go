package errors

import (
	"fmt"
)

// Raised when config file is missing
type ConfigFileMissingError string

func (str ConfigFileMissingError) Error() string {
	return fmt.Sprintf("ConfigFileMissingError: %q", string(str))
}

// Raised when the config file doesnt match schema
type ConfigValidationError string

func (str ConfigValidationError) Error() string {
	return fmt.Sprintf("ConfigValidationError: %q", string(str))
}

// Raised when the
type TemplateFileExistsError string

func (str TemplateFileExistsError) Error() string {
	return fmt.Sprintf("TemplateFileExistsError: %q", string(str))
}

// Raised when Question does not exist
type QuestionKeyInvalidError string

func (str QuestionKeyInvalidError) Error() string {
	return fmt.Sprintf("QuestionKeyInvalidError: %q", string(str))
}

type QuestionValidationError string

func (str QuestionValidationError) Error() string {
	return fmt.Sprintf("QuestionValidationError: %q", string(str))
}

// Raised when we cant convert the answer to the correct format (string, boolean, integer, etc)
type AnswerFormatError string

func (str AnswerFormatError) Error() string {
	return fmt.Sprintf("AnswerFormatError: %q", string(str))
}
