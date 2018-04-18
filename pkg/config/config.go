package config

import (
	"drawbridge/pkg/utils"
	"drawbridge/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"fmt"
	"bytes"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

// When initializing this class the following methods must be called:
// Config.New
// Config.Init
// This is done automatically when created via the Factory.
type configuration struct {
	*viper.Viper
}

//Viper uses the following precedence order. Each item takes precedence over the item below it:
// explicit call to Set
// flag
// env
// config
// key/value store
// default

func (c *configuration) Init() error {
	c.Viper = viper.New()
	//set defaults
	c.SetDefault("options.config_dir", "~/.ssh/drawbridge")
	c.SetDefault("options.pem_dir", "~/.ssh")
	c.SetDefault("options.active_config_template", "default")
	c.SetDefault("options.active_extra_templates", []string{})
	c.SetDefault("options.ui_group_priority", []string{"environment", "username"})
	c.SetDefault("options.ui_question_hidden", []string{})
	//TODO: options.overwrite == false/

	c.SetDefault("questions", map[string]Question{
		"pem_filename": {
			Description: "Pem key used to ssh to bastion",
			DefaultValue: "id_rsa",
			Schema: map[string]interface{}{
				"type": "string",
				"required": true,
			},
		},
		"environment": {
			Description: "Environment name for this stack",
			DefaultValue: "test",
			Schema: map[string]interface{}{
				"type": "string",
				"required": true,
			},
		},
		"username": {
			Description: "Username used to log into this stack",
			DefaultValue: "root",
			Schema: map[string]interface{}{
				"type": "string",
				"required": true,
			},
		},
		"domain": {
			Description: "Base domain name for all stacks",
			DefaultValue: "example.com",
			Schema: map[string]interface{}{
				"type": "string",
				"required": true,
			},
		},

	})
	c.SetDefault("answers", []map[string]interface{}{})

	c.SetDefault("config_templates.default.filepath", "{{.environment}}-{{.username}}-config")
	c.SetDefault("config_templates.default.content", utils.StripIndent(
	`
	ForwardAgent yes
	ForwardX11 no
	HashKnownHosts yes
	IdentitiesOnly yes
	StrictHostKeyChecking no

	Host bastion
	    Hostname bastion.example.com
	    User {{.username}}
	    IdentityFile {{.pem_dir}}/{{.pem_filename}}
	    LocalForward localhost:{{uniquePort .}} localhost:8080
	    UserKnownHostsFile=/dev/null
	    StrictHostKeyChecking=no
	`))


	//if you want to load a non-standard location system config file (~/drawbridge.yml), use ReadConfig
	c.SetConfigType("yaml")
	//c.SetConfigName("drawbridge")
	//c.AddConfigPath("$HOME/")

	//we're going to load the config file manually, since we need to validate it.
	err := c.ReadConfig("~/drawbridge.yaml") // Find and read the config file
	if _, ok := err.(errors.ConfigFileMissingError); ok { // Handle errors reading the config file
		//ignore "could not find config file"
	} else if (err != nil) {
		return err
	}

	//CLI options will be added via the `Set()` function

	return c.ValidateConfig()
}
func (c *configuration) ReadConfig(configFilePath string) error {
	configFilePath, err := utils.ExpandPath(configFilePath)
	if err != nil {
		return err
	}

	if !utils.FileExists(configFilePath) {
		log.Printf("No configuration file found at %v. Skipping", configFilePath)
		return errors.ConfigFileMissingError("The configuration file could not be found.")
	}

	//validate config file contents
	err = c.ValidateConfigFile(configFilePath)
	if err != nil {
		log.Printf("Config file at `%v` is invalid: %s", configFilePath, err)
		return err
	}

	log.Printf("Loading configuration file: %s", configFilePath)

	config_data, err := os.Open(configFilePath)
	if err != nil {
		log.Printf("Error reading configuration file: %s", err)
		return err
	}

	err = c.MergeConfig(config_data)
	if err != nil {
		return err
	}

	return c.ValidateConfig()
}

// This function ensures that the merged config works correctly.
func (c *configuration) ValidateConfig() error {

	////deserialize Questions
	//questionsMap := map[string]Question{}
	//err := c.UnmarshalKey("questions", &questionsMap)
	//
	//if err != nil {
	//	log.Printf("questions could not be deserialized correctly. %v", err)
	//	return err
	//}
	//
	//for _, v := range questionsMap {
	//
	//	typeContent, ok := v.Schema["type"].(string)
	//	if !ok || len(typeContent) == 0 {
	//		return errors.QuestionSyntaxError("`type` is required for questions")
	//	}
	//	//TODO: if a default value is set, check that the schema allows it.
	//}
	//
	////TODO: deserialize Answers
	//
	//// TODO: check if templates have any variables that are not defined as questions.

	return nil
}


func (c *configuration) ValidateConfigFile(configFilePath string) error {
	configFilePath, err := utils.ExpandPath(configFilePath)
	if err != nil {
		log.Printf("Could not expand filepath. %s", err)
		return err
	}

	configFileData, err := os.Open(configFilePath)
	if err != nil {
		log.Printf("Error reading configuration file: %s", err)
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(configFileData)
	configContent := map[string]interface{}{}
	err = yaml.Unmarshal(buf.Bytes(), &configContent)
	// To support boolean keys, the `yaml` package unmarshals maps to
	// map[interface{}]interface{}. Here we recurse through the result
	// and change all maps to map[string]interface{} like we would've
	// gotten from `json`.
	if err == nil {
		for k, v := range configContent {
			configContent[k] = utils.StringifyYAMLMapKeys(v)
		}
	} else {
		return err
	}

	// TODO: look at the dependencies key for matching the questions with answers keys.
	// TODO: look at the dependenices key for matching the options.active_templates with templates keys
	// TODO: ensure that all config_template.filepaths are relative, they will be created in the options.config_dir
	configFileSchema := `
	{
		"type": "object",
		"additionalProperties":false,
		"properties":{
			"options": {
				"type": "object",
				"additionalProperties":false,
				"properties": {
					"config_dir": {
						"type": "string"
					},
					"pem_dir": {
						"type": "string"
					},
					"active_config_template": {
						"type":"string"
					},
					"active_extra_templates": {
						"type":"array",
						"uniqueItems": true,
						"items":[{"type":"string"}]
					},
					"ui_group_priority": {
						"type":"array",
						"uniqueItems": true,
						"items":[{"type":"string"}],
						"maxItems": 3
					},
					"ui_question_hidden": {
						"type":"array",
						"uniqueItems": true,
						"items":[{"type":"string"}]
					}
				}
			},
			"questions":{
				"type": "object",
				"patternProperties": {
					"^[a-z0-9]*$":{
						"type":"object",
						"additionalProperties":false,
						"required": ["schema","description"],
						"properties": {
							"description": {
								"type": "string"
							},
							"default_value": {},
							"schema": {
								"type": "object",
								"additionalProperties":false,
								"required": ["type"],
								"properties": {
									"anyOf": {},
									"enum": {},
									"format": {},
									"maxLength": {},
									"maximum": {},
									"minLength": {},
									"minimum": {},
									"multipleOf": {},
									"not": {},
									"oneOf": {},
									"pattern": {},
									"required": {},
									"type": {
										"type": "string",
										"enum": ["integer", "number", "string", "boolean", "null"]
									}
								}
							}
						}
					}
				}
			},
			"answers":{
				"type": "array",
				"additionalProperties":false
			},
			"config_templates":{
				"type": "object",
				"additionalProperties":false
			},
			"extra_templates":{
				"type": "object",
				"additionalProperties":false
			}
		}
	}
	`

	schemaLoader := gojsonschema.NewStringLoader(configFileSchema)
	documentLoader := gojsonschema.NewGoLoader(configContent)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil{
		return err
	}
	if(!result.Valid()){
		errorMsg := ""
		for _, err := range result.Errors() {
			// Err implements the ResultError interface
			errorMsg += fmt.Sprintf("- %s\n", err)
		}

		return errors.QuestionValidationError(fmt.Sprintf("There was an error validating this config:\n %v ", errorMsg))
	}
	return nil

}

func (c *configuration) GetQuestion(questionKey string) (Question, error) {
	//deserialize Questions
	questions, err := c.GetQuestions()
	if err != nil {
		return Question{}, err
	}

	if question, ok := questions[questionKey]; ok {
		return question, nil
	} else {
		// the question does not exist
		return Question{}, errors.QuestionKeyInvalidError(fmt.Sprintf("There is no question for %v", questionKey))
	}
}

func (c *configuration) GetQuestions() (map[string]Question, error) {
	//deserialize Questions
	questionsMap := map[string]Question{}
	err := c.UnmarshalKey("questions", &questionsMap)
	return questionsMap, err
}

func (c *configuration) GetConfigTemplates() (map[string]Template, error) {
	//deserialize Templates
	templateMap := map[string]Template{}
	err := c.UnmarshalKey("config_templates", &templateMap)
	return templateMap, err
}

func (c *configuration) GetActiveConfigTemplate() (Template, error) {
	//deserialize Templates
	activeTemplateName := c.GetString("options.active_config_template")

	allTemplates, err := c.GetConfigTemplates()
	if err != nil{
		return Template{}, err
	}

	activeTemplate := allTemplates[activeTemplateName]
	return activeTemplate, nil
}


func (c *configuration) GetExtraTemplates() (map[string]Template, error) {
	//deserialize Templates
	templateMap := map[string]Template{}
	err := c.UnmarshalKey("extra_templates", &templateMap)
	return templateMap, err
}

func (c *configuration) GetActiveExtraTemplates() ([]Template, error) {
	//deserialize Templates
	activeTemplateNames := c.GetStringSlice("options.active_extra_templates")


	allTemplates, err := c.GetExtraTemplates()
	if err != nil{
		return nil, err
	}
	activeTemplates := []Template{}


	for _, activeTemplateName := range activeTemplateNames {
		activeTemplate := allTemplates[activeTemplateName]
		activeTemplates = append(activeTemplates, activeTemplate)
	}
	return activeTemplates, nil
}


