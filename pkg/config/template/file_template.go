package template

import (
	"drawbridge/pkg/errors"
	"drawbridge/pkg/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type FileTemplate struct {
	Template `mapstructure:",squash"`
	FilePath string `mapstructure:"filepath"`
}

func (t *FileTemplate) WriteTemplate(answerData map[string]interface{}) error {
	templatedFilePath, err := utils.PopulateTemplate(t.FilePath, answerData)
	if err != nil {
		return err
	}
	templatedFilePath, err = utils.ExpandPath(templatedFilePath)
	if err != nil {
		return err
	}
	answerData["filepath"] = templatedFilePath

	templatedContent, err := utils.PopulateTemplate(t.Content, answerData)
	if err != nil {
		return err
	}

	if !utils.FileExists(templatedFilePath) {

		//make the file's parent directory.
		err = os.MkdirAll(filepath.Dir(templatedFilePath), 0777)
		if err != nil {
			return err
		}

		log.Printf("Writing template to %v", templatedFilePath)
		err = utils.FileWrite(templatedFilePath, templatedContent, 0644)
		if err != nil {
			return err
		}

	} else {
		return errors.TemplateFileExistsError(fmt.Sprintf("file at %v already exists. Cannot write template file", templatedFilePath))
	}
	return nil
}
