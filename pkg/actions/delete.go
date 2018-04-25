package actions

import (
	"drawbridge/pkg/config"
	"drawbridge/pkg/utils"
	"fmt"
	"github.com/fatih/color"
	"path"
	"strings"
)

type DeleteAction struct {
	Config config.Interface
}

func (e *DeleteAction) All(answerDataList []map[string]interface{}, force bool) error {

	for _, v := range answerDataList {
		err := e.One(v, force)
		if err != nil {
			color.Red("ERROR IGNORED: %v", err)
		}
	}
	return nil
}
func (e *DeleteAction) One(answerData map[string]interface{}, force bool) error {

	//delete the config file by answerData
	renderedConfigFilePath := answerData["config"].(map[string]interface{})["filepath"].(string)

	//custom files specified in answerData
	renderedCustomFilePaths := []interface{}{}
	if customItems, ok := answerData["custom"]; ok && customItems != nil && len(customItems.([]interface{})) > 0 {
		renderedCustomFilePaths = customItems.([]interface{})
	}

	if !force {

		questionStr := []string{"Are you sure you would like to delete this config and associated templates? (PEM files will not be deleted)\n"}

		for k, v := range answerData {
			if utils.SliceIncludes(e.Config.InternalQuestionKeys(), k) {
				continue
			}
			questionStr = append(questionStr, fmt.Sprintf("%v: %v", k, v))
		}
		questionStr = append(questionStr, "\nPlease confirm [yes/no]:")

		val := utils.StdinQueryBoolean(strings.Join(questionStr, "\n"))
		if !val {
			color.Red("Cancelled delete operation.")
			return nil
		}
	}

	fmt.Printf("Deleting config file: %v\n", renderedConfigFilePath)
	if utils.FileExists(renderedConfigFilePath) {
		utils.FileDelete(renderedConfigFilePath)
	} else {
		color.Yellow(" - Skipping. Could not find config file at: %v", renderedConfigFilePath)
	}

	//delete any custom templates.
	for _, customTemplateData := range renderedCustomFilePaths {
		renderedCustomFilePath := customTemplateData.(map[string]interface{})["filepath"].(string)
		fmt.Printf("Deleting custom file: %v\n", renderedCustomFilePath)
		if utils.FileExists(renderedCustomFilePath) {
			utils.FileDelete(renderedCustomFilePath)
		} else {
			color.Yellow(" - Skipping. Could not find config file at: %v", renderedCustomFilePath)
		}
	}
	//delete the .answers.yaml
	fmt.Println("Deleting answers file")
	answersFilePath := path.Join(answerData["config_dir"].(string), fmt.Sprintf(".%v.answers.yaml", path.Base(renderedConfigFilePath)))
	if utils.FileExists(answersFilePath) {
		utils.FileDelete(answersFilePath)
	} else {
		color.Yellow(" - Skipping. Could not find answers file at: %v", answersFilePath)
	}

	return nil
}
