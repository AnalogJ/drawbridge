package actions

import (
	"drawbridge/pkg/config"
	//"io/ioutil"
	"bytes"
	"drawbridge/pkg/utils"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ListAction struct {
	Config         config.Interface
	GroupedAnswers *gabs.Container
	OrderedAnswers []interface{}
}

func (e *ListAction) Start() error {

	answersList, err := e.RenderedAnswersList()
	if err != nil {
		return err
	}

	e.GroupedAnswers = e.GroupAnswerList(answersList, e.Config.GetStringSlice("options.ui_group_priority"))

	e.PrintUI(e.GroupedAnswers)

	return nil
}

//this is a list of all the answers that have been used to populate templates & config files.
//will be ordered by config file name
func (e *ListAction) RenderedAnswersList()([]map[string]interface{}, error) {
	files, err := e.FindAllAnswerFiles()
	if err != nil {
		return nil, err
	}

	return e.ParseAnswerFiles(files)
}

func (e *ListAction) FindAllAnswerFiles() ([]string,error) {
	configDir := e.Config.GetString("options.config_dir")
	configDir, err := utils.ExpandPath(configDir)
	if err != nil {
		return nil, err
	}
	// files, err := ioutil.ReadDir(configDir)
	return filepath.Glob(filepath.Join(configDir, ".*.answers.yaml"))
}

func (e *ListAction) ParseAnswerFiles(files []string) ([]map[string]interface{}, error) {
	answersList := []map[string]interface{}{}
	for _, answerFilePath := range files {

		//read file
		answerFileData, err := os.Open(answerFilePath)
		if err != nil {
			log.Printf("Error reading answer file: %s", err)
			return nil, err
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
			return nil, err
		}

		for k, v := range answerData {
			answerData[k] = utils.StringifyYAMLMapKeys(v)
		}

		answersList = append(answersList, answerData)
	}
	return answersList, nil
}

func (e *ListAction) GroupAnswerList(answersList []map[string]interface{}, groupKeys []string) (*gabs.Container) {
	// Group By for existing configs.

	groupedAnswers := gabs.New()
	if len(groupKeys) > 0 {

		for _, answerData := range answersList {
			keyValues := []string{}
			for _, questionKey := range groupKeys {
				if value, ok := answerData[questionKey]; ok {
					keyValues = append(keyValues, fmt.Sprintf("%v", value))
				} else {
					keyValues = append(keyValues, "")
				}
			}

			// now make sure we have an array at this level.
			if !groupedAnswers.Exists(keyValues...) {
				groupedAnswers.Array(keyValues...)
			}
			groupedAnswers.ArrayAppend(answerData, keyValues...)
		}

	} else {
		groupedAnswers.Set(answersList, "")
	}
	return groupedAnswers
}


func (e *ListAction) PrintUI(groupedAnswers *gabs.Container) error {
	return e.recursivePrintUI(0, []string{}, groupedAnswers)
}


func (e *ListAction) recursivePrintUI(level int, groups []string, groupedAnswers *gabs.Container) error {
	children, _ := groupedAnswers.ChildrenMap()
	for groupKey, child := range children {
		nextGroups := []string{}
		if level == 0 {
			coloredPrintf(level, "%v\n", groupKey)
		} else {
			nextGroups = append(nextGroups, groups...)
			nextGroups = append(nextGroups, groupKey)
		}

		switch v := child.Data().(type) {
		case map[string]interface{}:
			e.recursivePrintUI(level+1, nextGroups, child)
		case []interface{}:

			printGroupHeader(nextGroups)

			for _, answer := range child.Data().([]interface{}) {
				e.OrderedAnswers = append(e.OrderedAnswers, answer)

				answerStr := printAnswer(len(e.OrderedAnswers), answer.(map[string]interface{}), e.Config.GetStringSlice("options.ui_question_hidden"), e.Config.GetStringSlice("options.ui_group_priority"))
				fmt.Printf(answerStr)

			}
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}
	}
	return nil

}

func printGroupHeader(secondaryGroups []string) {
	header := ":::: "
	if len(secondaryGroups) >= 1 {
		header += fmt.Sprintf("%v ", color.GreenString(secondaryGroups[0]))
	}
	if len(secondaryGroups) >= 2 {
		header += fmt.Sprintf("%v ", color.CyanString(secondaryGroups[1]))
	}
	if len(secondaryGroups) >= 3 {
		header += fmt.Sprintf("(%v) ", color.YellowString(secondaryGroups[2]))
	}

	maxLength := 50

	header += fmt.Sprintf("%v\n", strings.Repeat(":", (maxLength-(1+len(header)))))

	fmt.Print(header)
}

func printAnswer(id int, answer map[string]interface{}, uiHiddenKeys []string, uiGroupPriority []string) string {

	//fmt.Printf("\t%v\t%v\n", color.YellowString(strconv.Itoa(), answerStr)

	answerStr := fmt.Sprintf("\t%v\t", color.YellowString(strconv.Itoa(id)))
	for k, v := range answer {
		if utils.StringInSlice(uiHiddenKeys, k) || utils.StringInSlice(uiGroupPriority, k) {
			continue
		}
		answerStr += fmt.Sprintf("%v: %v\n\t\t", k, v)
	}
	answerStr += "\n"
	return answerStr
}

func coloredPrintf(level int, formattedStr string, data ...interface{}) {
	if level == 0 {
		color.Red(formattedStr, data...)
	} else if level == 1 {
		color.Green(formattedStr, data...)
	} else if level == 2 {
		color.Cyan(formattedStr, data...)
	} else {
		fmt.Print("Unkonw int type")
	}
}
