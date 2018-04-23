package template

import (
	"drawbridge/pkg/utils"
	"github.com/fatih/color"
	"log"
	"os"
	"path/filepath"
)

type PacTemplate struct {
	FileTemplate `mapstructure:",squash"`
}

func (t *PacTemplate) WriteTemplate(answerDataList []map[string]interface{}, dryRun bool) (map[string]interface{}, error) {
	if t.data == nil {
		t.data = map[string]interface{}{}
	}

	pacFilePath, err := utils.ExpandPath(t.FilePath)
	if err != nil {
		return nil, err
	}

	t.data["filepath"] = pacFilePath

	templatedContent, err := utils.PopulateTemplate(t.Content, answerDataList)
	if err != nil {
		return nil, err
	}

	if !utils.FileExists(pacFilePath) {

		//make the file's parent directory.
		err = os.MkdirAll(filepath.Dir(pacFilePath), 0777)
		if err != nil {
			return nil, err
		}
	} else {
		color.Yellow("Pac file already exists, updating.")
	}
	log.Printf("Writing template to %v", pacFilePath)
	err = utils.FileWrite(pacFilePath, templatedContent, 0644, dryRun)
	if err != nil {
		return nil, err
	}

	return t.data, nil
}
