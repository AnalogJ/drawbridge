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

// Raised when a dependency (like ssh or ssh-agent) is missing
type DependencyMissingError string

func (str DependencyMissingError) Error() string {
	return fmt.Sprintf("DependencyMissingError: %q", string(str))
}

// Raised when a SSH/pem key is not present.
type PemKeyMissingError string

func (str PemKeyMissingError) Error() string {
	return fmt.Sprintf("PemKeyMissingError: %q", string(str))
}

// Raised when the file to write already exists
type TemplateFileExistsError string

func (str TemplateFileExistsError) Error() string {
	return fmt.Sprintf("TemplateFileExistsError: %q", string(str))
}

// Raised when Question does not exist
type QuestionKeyInvalidError string

func (str QuestionKeyInvalidError) Error() string {
	return fmt.Sprintf("QuestionKeyInvalidError: %q", string(str))
}

type AnswerValidationError string

func (str AnswerValidationError) Error() string {
	return fmt.Sprintf("AnswerValidationError: %q", string(str))
}

// Raised when we cant convert the answer to the correct format (string, boolean, integer, etc)
type AnswerFormatError string

func (str AnswerFormatError) Error() string {
	return fmt.Sprintf("AnswerFormatError: %q", string(str))
}
