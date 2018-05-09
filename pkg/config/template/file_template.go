package template

import (
	"drawbridge/pkg/errors"
	"drawbridge/pkg/utils"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"path/filepath"
)

type FileTemplate struct {
	Template `mapstructure:",squash"`
	FilePath string `mapstructure:"filepath"`
}

//func (t *FileTemplate) PopulateFilePath(answerData map[string]interface{}) (string, error) {
//	templatedFilePath, err := utils.PopulateTemplate(t.FilePath, answerData)
//	if err != nil {
//		return "", err
//	}
//	templatedFilePath, err = utils.ExpandPath(templatedFilePath)
//	if err != nil {
//		return "", err
//	}
//	return templatedFilePath, nil
//}

func (t *FileTemplate) DeleteTemplate(answerData map[string]interface{}) error {
	templatedFilePath, err := utils.PopulatePathTemplate(t.FilePath, answerData)
	if err != nil {
		return nil
	}

	if !utils.FileExists(templatedFilePath) {
		// warn that this file does not exist
		color.Yellow(" - Skipping. Could not find file: %v", templatedFilePath)
		return nil
	} else {
		return os.Remove(templatedFilePath)
	}
}

func (t *FileTemplate) WriteTemplate(answerData map[string]interface{}, dryRun bool) (map[string]interface{}, error) {
	if t.data == nil {
		t.data = map[string]interface{}{}
	}

	answerData, err := utils.MapDeepCopy(answerData)
	if err != nil {
		return nil, err
	}

	templatedFilePath, err := utils.PopulatePathTemplate(t.FilePath, answerData)
	if err != nil {
		return nil, err
	}

	t.data["filepath"] = templatedFilePath
	answerData["template"] = t.data

	templatedContent, err := utils.PopulateTemplate(t.Content, answerData)
	if err != nil {
		return nil, err
	}

	if !utils.FileExists(templatedFilePath) {

		//make the file's parent directory.
		err = os.MkdirAll(filepath.Dir(templatedFilePath), 0777)
		if err != nil {
			return nil, err
		}

		log.Printf("Writing template to %v", templatedFilePath)
		err = utils.FileWrite(templatedFilePath, templatedContent, 0644, dryRun)
		if err != nil {
			return nil, err
		}

	} else {
		return nil, errors.TemplateFileExistsError(fmt.Sprintf("file at %v already exists. Cannot write template file", templatedFilePath))
	}

	return t.data, nil
}
