package config

import (
	"fmt"
	"github.com/analogj/drawbridge/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

type Question struct {
	Description  string                 `mapstructure:"description"`
	DefaultValue interface{}            `mapstructure:"default_value"`
	Schema       map[string]interface{} `mapstructure:"schema"`
}

func (q *Question) GetType() string {
	return q.Schema["type"].(string)
}

func (q *Question) Required() bool {
	isRequired, isSet := q.Schema["required"].(bool)
	return isRequired && isSet
}

func (q *Question) Validate(questionKey string, answerValue interface{}) error {
	questionSchema := map[string]interface{}{
		"properties": map[string]map[string]interface{}{
			questionKey: {},
		},
		"required": []string{},
	}

	isRequired := q.Schema["required"]
	if isRequired != nil {
		questionSchema["required"] = append(questionSchema["required"].([]string), questionKey)
	}

	//fix viper case-insensitivity & cleanup Schema
	properRuleKeys := map[string]string{
		"allof":             "allOf",
		"anyof":             "anyOf",
		"maxitems":          "maxItems",
		"maxlength":         "maxLength",
		"maxproperties":     "maxProperties",
		"minitems":          "minItems",
		"minlength":         "minLength",
		"minproperties":     "minProperties",
		"multipleof":        "multipleOf",
		"oneof":             "oneOf",
		"patternproperties": "patternProperties",
		"uniqueitems":       "uniqueItems",
	}

	for ruleKey, ruleValue := range q.Schema {
		if ruleKey == "required" {
			//skip, required is already handled above.
			continue
		}

		actualKey := ""
		if val, ok := properRuleKeys[ruleKey]; ok {
			//lets fix the rule key to use the uppercase version.
			//fmt.Printf("\nSwitching %v for %v\n", actualKey, val)
			actualKey = val
		} else {
			actualKey = ruleKey
		}

		questionSchema["properties"].(map[string]map[string]interface{})[questionKey][actualKey] = ruleValue
	}

	schemaLoader := gojsonschema.NewGoLoader(questionSchema)

	testData := map[string]interface{}{
		questionKey: answerValue,
	}

	documentLoader := gojsonschema.NewGoLoader(testData)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		//TODO: populate with actual errors from result obj.

		return errors.AnswerValidationError(fmt.Sprintf("There was an error validating this answer: %v", result.Errors()))
	}
	return nil
}
