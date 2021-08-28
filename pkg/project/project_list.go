package project

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/analogj/drawbridge/pkg/errors"
	"github.com/analogj/drawbridge/pkg/utils"
	"github.com/fatih/color"
	"github.com/xlab/treeprint"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type ProjectList struct {
	projects []projectData

	groupByKeys []string
	hiddenKeys  []string

	groupedAnswers     *gabs.Container
	groupedAnswersList []map[string]interface{}
	groupedTree        treeprint.Tree
}

// PUBLIC functions

func (p *ProjectList) Length() int {
	return len(p.projects)
}

func (p *ProjectList) GetAll() []map[string]interface{} {
	if p.Length() == 0 {
		return []map[string]interface{}{}
	}

	if len(p.groupedAnswersList) == 0 {
		p.initGroups()
	}

	return p.groupedAnswersList
}

func (p *ProjectList) GetIndex(index_0based int) (map[string]interface{}, error) {
	if p.Length() == 0 {
		return nil, errors.ProjectListEmptyError("No answers found, please call `drawbridge create` first")
	}

	if len(p.groupedAnswersList) == 0 {
		p.initGroups()
	}

	if index_0based < 0 || index_0based >= len(p.projects) {
		return nil, errors.ProjectListIndexInvalidError(fmt.Sprintf("Selected index (%v) is invalid. Must be between %v-%v", index_0based+1, 1, p.Length()))
	} else {
		return p.groupedAnswersList[index_0based], nil
	}
}

func (p *ProjectList) SetAliasForIndex(index_0based int, alias string) (map[string]interface{}, error) {

	answerYamlFile, err := ioutil.ReadFile(p.projects[index_0based].AnswerFilePath)
	if err != nil {
		return nil, errors.ConfigFileMissingError("could not open answerfile for config")
	}

	answerData := make(map[string]interface{})

	err = yaml.Unmarshal(answerYamlFile, &answerData)
	if err != nil {
		return nil, errors.ConfigFileMissingError("could not parse answerfile")
	}

	if existingAlias, existingAliasOk := answerData["alias"]; existingAliasOk {
		color.HiYellow("Warning: replacing existing alias (%s) with new value: %s", existingAlias, alias)
	}
	answerData["alias"] = alias

	answersFileContent, err := yaml.Marshal(answerData)
	if err != nil {
		return nil, err
	}
	err = utils.FileWrite(p.projects[index_0based].AnswerFilePath, string(answersFileContent), 0640, false)

	return answerData, err
}

func (p *ProjectList) Prompt(message string) (map[string]interface{}, int, error) {
	if p.Length() == 0 {
		return nil, 0, errors.ProjectListEmptyError("No answers found, please call `drawbridge create` first")
	}

	if len(p.groupedAnswersList) == 0 {
		p.initGroups()
	}

	p.PrintTree("")

	for true {

		//prompt the user to enter a valid choice
		index_1based, err := utils.StdinQueryInt(fmt.Sprintf("%v (%v-%v):", message, 1, p.Length()))
		if err != nil {
			color.HiRed("ERROR: %v", err)
			continue
		}

		if !(index_1based > 0 && index_1based <= p.Length()) {
			color.HiRed("Invalid selection. Must be between %v-%v", 1, p.Length())
			continue
		}

		return p.groupedAnswersList[index_1based-1], index_1based - 1, nil
	}
	return nil, 0, nil
}

func (p *ProjectList) PrintTree(startMessage string) {
	p.groupedTree.SetValue("Rendered Drawbridge Configs:")
	fmt.Println(p.groupedTree.String())
}

// Private functions

func (p *ProjectList) initGroups() {
	//intialize storage
	p.groupedAnswers = gabs.New()
	p.groupedAnswersList = []map[string]interface{}{}
	p.groupedTree = treeprint.New()

	//group the project answers
	p.groupProjectAnswers()
	//populate the ordered group list, and tree
	p.recursivePopulateGroupListAndTree(0, p.groupedTree, p.groupedAnswers)

}

func (p *ProjectList) groupProjectAnswers() {
	// Group By for existing configs.

	if len(p.groupByKeys) > 0 {

		for _, project := range p.projects {
			keyValues := []string{}
			for _, questionKey := range p.groupByKeys {
				if value, ok := project.Answers[questionKey]; ok && value != nil {
					keyValues = append(keyValues, fmt.Sprintf("%v", value))
				} else {
					keyValues = append(keyValues, "")
				}
			}

			// now make sure we have an array at this level.
			if !p.groupedAnswers.Exists(keyValues...) {
				p.groupedAnswers.Array(keyValues...)
			}
			p.groupedAnswers.ArrayAppend(project.Answers, keyValues...)
		}

	} else {

		answersList := []map[string]interface{}{}
		for _, project := range p.projects {
			answersList = append(answersList, project.Answers)
		}

		p.groupedAnswers.Set(answersList, "")
	}
}

func (p *ProjectList) recursivePopulateGroupListAndTree(level int, parentTree treeprint.Tree, groupedAnswers *gabs.Container) error {

	groupByKeys := p.groupByKeys

	children, _ := groupedAnswers.ChildrenMap()

	groupKeys := []string{}
	for k := range children {
		groupKeys = append(groupKeys, k)
	}
	sort.Strings(groupKeys)

	for _, groupKey := range groupKeys {
		child := children[groupKey]
		currentTree := parentTree

		//ensure the current groupKey is not empty.
		if len(groupKey) > 0 {

			// handle following cases:
			if level+1 < len(groupByKeys) {
				currentTree = parentTree.AddMetaBranch(p.coloredString(level, groupKey), groupByKeys[level])
			}
		}

		switch v := child.Data().(type) {
		case map[string]interface{}:
			p.recursivePopulateGroupListAndTree(level+1, currentTree, child)
		case []interface{}:

			//printGroupHeader(nextGroups)

			answerList := child.Data().([]interface{})
			sort.Slice(answerList, func(i, j int) bool {
				iItem := answerList[i].(map[string]interface{})
				jItem := answerList[j].(map[string]interface{})

				if iItem[groupKey] != nil && jItem[groupKey] != nil {
					return iItem[groupKey].(string) > jItem[groupKey].(string)
				} else {
					return false
				}
			})

			for _, answer := range answerList {
				p.groupedAnswersList = append(p.groupedAnswersList, answer.(map[string]interface{}))

				//answerStr := printAnswer(len(e.OrderedAnswers), answer.(map[string]interface{}), e.Config.GetStringSlice("options.ui_question_hidden"), e.Config.GetStringSlice("options.ui_group_priority"))
				alias, aliasOk := answer.(map[string]interface{})["alias"]
				answerIndex := strconv.Itoa(len(p.groupedAnswersList))
				if aliasOk {
					answerIndex = fmt.Sprintf("%s, %s", answerIndex, alias)
				}

				currentTree.AddMetaNode(
					color.YellowString(answerIndex),
					p.answerString(groupByKeys[level], answer.(map[string]interface{})))
			}
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}
	}
	return nil
}

func (p *ProjectList) answerString(highlightGroupKey string, answer map[string]interface{}) string {

	answerStr := []string{color.BlueString(fmt.Sprintf("%v: %v", highlightGroupKey, answer[highlightGroupKey]))}

	keys := utils.MapKeys(answer)

	for _, k := range keys {
		v := answer[k]

		//skip hidden keys, group by keys and internal keys.
		if utils.SliceIncludes(p.hiddenKeys, k) || utils.SliceIncludes(p.groupByKeys, k) {
			continue
		}

		//skip highlighted group
		if k == highlightGroupKey {
			continue
		}

		answerStr = append(answerStr, fmt.Sprintf("%v: %v", k, v))
	}
	return strings.Join(answerStr, ", ")
}

func (p *ProjectList) coloredString(level int, data string) string {
	if level == 0 {
		return color.RedString(data)
	} else if level == 1 {
		return color.GreenString(data)
	} else if level == 2 {
		return color.CyanString(data)
	} else {
		return data
	}
}
