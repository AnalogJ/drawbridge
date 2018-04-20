package list

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

type ListEngine struct {
	Config         config.Interface
	GroupedAnswers *gabs.Container
	OrderedAnswers []interface{}
}

func (e *ListEngine) Start() error {

	configDir := e.Config.GetString("options.config_dir")
	configDir, err := utils.ExpandPath(configDir)
	if err != nil {
		return err
	}
	// files, err := ioutil.ReadDir(configDir)
	files, err := filepath.Glob(filepath.Join(configDir, ".*.answers.yaml"))
	if err != nil {
		return err
	}

	answersList := []map[string]interface{}{}
	for _, answerFilePath := range files {
		fmt.Println(answerFilePath)

		//read file
		answerFileData, err := os.Open(answerFilePath)
		if err != nil {
			log.Printf("Error reading answer file: %s", err)
			return err
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
			return err
		}

		for k, v := range answerData {
			answerData[k] = utils.StringifyYAMLMapKeys(v)
		}

		answersList = append(answersList, answerData)
	}
	fmt.Printf("\nANSWERS LIST\n%v", answersList)

	// Group By for existing configs.
	priorityOrder := e.Config.GetStringSlice("options.ui_group_priority")
	groupedAnswers := gabs.New()
	if len(priorityOrder) > 0 {

		for _, answerData := range answersList {
			//subGroup := groupedAnswers
			//for index, questionKey := range priorityOrder {
			//	if _, ok := subGroup[questionKey].(map[string]interface{}); !ok {
			//		// this level of grouping doesnt exist, we need to create it.
			//		if(index == len(priorityOrder)){
			//			subGroup[questionKey] = []map[string]interface{}{};
			//		} else {
			//			subGroup[questionKey] = map[string]interface{}{};
			//		}
			//	}
			//	subGroup = subGroup[questionKey]
			//}
			//subGroup = append(subGroup, answerData)

			keyValues := []string{}
			for _, questionKey := range priorityOrder {
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
		//groupedAnswers.Array("")
		groupedAnswers.Set(answersList, "")

	}
	e.GroupedAnswers = groupedAnswers
	fmt.Printf("\nGROUPED ANSWERSLIST\n%v", groupedAnswers)

	e.PrintUI(0, []string{}, groupedAnswers)

	return nil
}

func (e *ListEngine) PrintUI(level int, groups []string, groupedAnswers *gabs.Container) error {
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
			e.PrintUI(level+1, nextGroups, child)
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
