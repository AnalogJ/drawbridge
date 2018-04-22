package actions

import (
	"drawbridge/pkg/config"
	"drawbridge/pkg/errors"
	"drawbridge/pkg/utils"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"path"
	"sort"
	"strconv"
)

type CreateAction struct {
	Config config.Interface
}

func (e *CreateAction) Start(cliAnswerData map[string]interface{}) error {

	// prepare answer data with config.options
	answerData := map[string]interface{}{}
	e.Config.UnmarshalKey("options", &answerData)

	// add defaults into answerData
	questions, err := e.Config.GetQuestions()
	if err != nil {
		return err
	}
	for questionKey, question := range questions {
		if question.DefaultValue != nil {
			answerData[questionKey] = question.DefaultValue
		}
	}

	// merge cliAnswerData into answerData
	for cliAnswerKey, cliAnswerValue := range cliAnswerData {
		answerData[cliAnswerKey] = cliAnswerValue
	}

	//log.Printf("answers found before questioning: %v \n", answerData)

	fmt.Println("Current Answers:")

	questionKeys := utils.MapKeys(answerData)
	for _, questionKey := range questionKeys {
		if utils.StringInSlice(e.Config.InternalQuestionKeys(), questionKey) {
			continue
		}

		fmt.Printf("%v: %v\n",
			questionKey,
			color.GreenString(fmt.Sprintf("%v", answerData[questionKey])))
	}

	// ensuer that that all questions are answered, query user if missing anything.
	answerData, err = e.Query(questions, answerData)
	if err != nil {
		return err
	}

	//set any optional keys to nil value.
	for questionKey, question := range questions {
		if !question.Required() {

			if _, ok := answerData[questionKey]; !ok {
				//answerdata does not contain this optional key
				answerData[questionKey] = nil
			}
		}
	}

	// write the config template, make sure we "fix" the config filepath
	activeConfigTemplate, err := e.Config.GetActiveConfigTemplate()
	if err != nil {
		return err
	}

	err = activeConfigTemplate.WriteTemplate(answerData, e.Config.InternalQuestionKeys())
	if err != nil {
		return err
	}

	// load up all active_custom_templates and attempt to merge answers with it.
	activeCustomTemplates, err := e.Config.GetActiveCustomTemplates()
	if err != nil {
		return err
	}

	for _, template := range activeCustomTemplates {
		err := template.WriteTemplate(answerData)
		if err != nil {
			return err
		}
	}

	// write the answers.yaml file
	answersFilePath := path.Join(e.Config.GetString("options.config_dir"), fmt.Sprintf(".%v.answers.yaml", path.Base(activeConfigTemplate.FilePath)))
	answersFilePath, err = utils.PopulateTemplate(answersFilePath, answerData)
	if err != nil {
		return err
	}

	answersFileContent, err := yaml.Marshal(answerData)
	if err != nil {
		return err
	}
	err = utils.FileWrite(answersFilePath, string(answersFileContent), 0600)
	if err != nil {
		return err
	}

	return nil
}

func (e *CreateAction) Query(questions map[string]config.Question, answerData map[string]interface{}) (map[string]interface{}, error) {

	questionKeys := []string{}
	for k := range questions {
		questionKeys = append(questionKeys, k)
	}
	sort.Strings(questionKeys)

	for _, questionKey := range questionKeys {
		questionData := questions[questionKey]

		val, ok := questionData.Schema["required"]
		required := ok && val.(bool)

		if _, ok := answerData[questionKey]; !ok && required {
			answerData[questionKey] = e.queryResponse(questionKey, questionData)

		}
	}

	return answerData, nil
}

func (e *CreateAction) queryResponse(questionKey string, question config.Question) interface{} {

	for true {
		//this question is not answered, and it is required. We should ask the user.
		answer := utils.StdinQuery(fmt.Sprintf("Please enter a value for `%s` [%s] - %s:", questionKey, question.GetType(), question.Description))

		answerTyped, err := convertAnswerType(answer, question.GetType())
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
		//TODO: figure out how to handle empty strings (still valid answer for some reason)

		err = question.Validate(questionKey, answerTyped)
		if err != nil {
			color.HiRed("%v\n", err)
			//fmt.Printf("%v\n", err)
		} else {
			return answerTyped
		}

	}
	//return answerTyped
	return nil
}

func convertAnswerType(answer string, questionType string) (interface{}, error) {
	if questionType == "integer" {
		answer, err := strconv.ParseInt(answer, 10, 64)
		if err != nil {
			return nil, err
		}
		return answer, nil
	} else if questionType == "number" {
		answer, err := strconv.ParseFloat(answer, 64)
		if err != nil {
			return nil, err
		}
		return answer, nil
	} else if questionType == "boolean" {
		answer, err := strconv.ParseBool(answer)
		if err != nil {
			return nil, err
		}
		return answer, nil
	} else if questionType == "string" {
		return answer, nil
	} else {
		return nil, errors.AnswerFormatError(fmt.Sprintf("could not convert %v to unknown %v type", answer, questionType))
	}

}
