package project

import (
	"drawbridge/pkg/config"
	"drawbridge/pkg/utils"
	"path/filepath"
	"os"
	"log"
	"bytes"
	"gopkg.in/yaml.v2"
)

// this will populate the ProjectList with answers loaded from the filesystem
func CreateProjectListFromConfigDir(configData config.Interface) (ProjectList, error) {

	projectList := ProjectList{
		projects: []projectData{},
		groupByKeys: configData.GetStringSlice("options.ui_group_priority"),

	}
	projectList.hiddenKeys = append(projectList.hiddenKeys, configData.GetStringSlice("options.ui_question_hidden")...)
	projectList.hiddenKeys = append(projectList.hiddenKeys, configData.InternalQuestionKeys()...)


	answerFiles, err := answerFilesInConfigDir(configData.GetString("options.config_dir"))
	if err != nil{
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
func CreateProjectListFromProvidedAnswers(configData config.Interface)(ProjectList, error){
	projectList := ProjectList{
		projects: []projectData{},
		groupByKeys: configData.GetStringSlice("options.ui_group_priority"),
	}
	projectList.hiddenKeys = append(projectList.hiddenKeys, configData.GetStringSlice("options.ui_question_hidden")...)
	projectList.hiddenKeys = append(projectList.hiddenKeys, configData.InternalQuestionKeys()...)

	providedAnswerDataList, err := configData.GetProvidedAnswerList()
	if err != nil{
		return projectList, err
	}

	for _, answerData := range  providedAnswerDataList {
		providedProjectData := projectData{
			Answers: answerData,
		}
		projectList.projects = append(projectList.projects, providedProjectData)
	}

	return projectList, nil

}


func CreateProjectFromConfigDirAnswerFile(configFilePath string)(projectData, error){
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

	return projectData{
		Answers: answerData,
		AnswerFilePath: answerFilePath,
		ConfigFilePath: answerData["config"].(map[string]interface{})["filepath"].(string),
		PemFilePath: answerData["config"].(map[string]interface{})["pem_filepath"].(string),
	}, nil

}