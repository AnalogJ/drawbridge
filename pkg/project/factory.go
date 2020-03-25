package project

import (
	"bytes"
	"github.com/analogj/drawbridge/pkg/config"
	"github.com/analogj/drawbridge/pkg/utils"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

// this will populate the ProjectList with answers loaded from the filesystem
func CreateProjectListFromConfigDir(configData config.Interface) (ProjectList, error) {

	projectList := ProjectList{
		projects:    []projectData{},
		groupByKeys: configData.GetStringSlice("options.ui_group_priority"),
	}
	projectList.hiddenKeys = append(projectList.hiddenKeys, configData.GetStringSlice("options.ui_question_hidden")...)
	projectList.hiddenKeys = append(projectList.hiddenKeys, configData.InternalQuestionKeys()...)

	answerFiles, err := answerFilesInConfigDir(configData.GetString("options.config_dir"))
	if err != nil {
		return projectList, err
	}

	for _, answerFilePath := range answerFiles {

		answerProjectData, err := CreateProjectFromConfigDirAnswerFile(answerFilePath)
		if err != nil {
			//TODO do somethignt here with the error (print it out?)
			continue
		}
		projectList.projects = append(projectList.projects, answerProjectData)

	}

	return projectList, nil
}

//this will populate the ProjectList with answers emedded in the config file (~/drawbridge.yaml)
func CreateProjectListFromProvidedAnswers(configData config.Interface) (ProjectList, error) {
	projectList := ProjectList{
		projects:    []projectData{},
		groupByKeys: configData.GetStringSlice("options.ui_group_priority"),
	}
	projectList.hiddenKeys = append(projectList.hiddenKeys, configData.GetStringSlice("options.ui_question_hidden")...)
	projectList.hiddenKeys = append(projectList.hiddenKeys, configData.InternalQuestionKeys()...)

	providedAnswerDataList, err := configData.GetProvidedAnswerList()
	if err != nil {
		return projectList, err
	}

	for _, answerData := range providedAnswerDataList {
		providedProjectData := projectData{
			Answers: answerData,
		}
		projectList.projects = append(projectList.projects, providedProjectData)
	}

	return projectList, nil

}

func CreateProjectFromConfigDirAnswerFile(configFilePath string) (projectData, error) {
	return parseAnswerFile(configFilePath)
}

///////////////////////////////////////////////////////////////////////////////
// Helpers

func answerFilesInConfigDir(configDir string) ([]string, error) {
	configDir, err := utils.ExpandPath(configDir)
	if err != nil {
		return nil, err
	}
	// files, err := ioutil.ReadDir(configDir)
	return filepath.Glob(filepath.Join(configDir, ".*.answers.yaml"))
}

func parseAnswerFile(answerFilePath string) (projectData, error) {

	//read file
	answerFileData, err := os.Open(answerFilePath)
	if err != nil {
		log.Printf("Error reading answer file: %s", err)
		return projectData{}, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(answerFileData)
	answerData := map[string]interface{}{}
	err = yaml.Unmarshal(buf.Bytes(), &answerData)
	// To support boolean keys, the `yaml` package unmarshals maps to
	// map[interface{}]interface{}. Here we recurse through the result
	// and change all maps to map[string]interface{} like we would've
	// gotten from `json`.
	if err != nil {
		return projectData{}, err
	}

	for k, v := range answerData {
		answerData[k] = utils.StringifyYAMLMapKeys(v)
	}

	//TODO: warn the user if the answer data would no longer render the same answers.yaml file.

	answerDataConfig := answerData["config"].(map[string]interface{})
	pemFilePath := "" //this is an optional field (may be unset/nil in some configs)
	if val, ok := answerDataConfig["pem_filepath"]; ok {
		pemFilePath = val.(string)
	}

	return projectData{
		Answers:        answerData,
		AnswerFilePath: answerFilePath,
		ConfigFilePath: answerDataConfig["filepath"].(string),
		PemFilePath:    pemFilePath,
	}, nil

}
