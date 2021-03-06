package config

import (
	"bytes"
	"fmt"
	"github.com/analogj/drawbridge/pkg/config/template"
	"github.com/analogj/drawbridge/pkg/errors"
	"github.com/analogj/drawbridge/pkg/utils"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
	"log"
	"os"
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
	c.SetDefault("options.pem_dir", "~/.ssh/drawbridge/pem")
	c.SetDefault("options.active_config_template", "default")

	c.SetDefault("options.active_custom_templates", []string{})
	c.SetDefault("options.ui_group_priority", []string{"environment", "stack_name", "shard", "shard_type"})
	c.SetDefault("options.ui_question_hidden", []string{})

	c.SetDefault("questions", map[string]Question{
		"environment": {
			Description: "What is the environment name?",
			Schema: map[string]interface{}{
				"type":     "string",
				"required": true,
				"enum":     []string{"test", "stage", "prod"},
			},
		},
		"stack_name": {
			Description:  "What is the stack name?",
			DefaultValue: "app",
			Schema: map[string]interface{}{
				"type":      "string",
				"required":  true,
				"minLength": 1,
				"maxLength": 6,
			},
		},
		"shard": {
			Description: "What is the shard datacenter?",
			Schema: map[string]interface{}{
				"type":     "string",
				"required": true,
				"enum":     []string{"us-east-1", "us-east-2", "eu-west-1", "eu-west-2", "ap-south-1"},
			},
		},
		"shard_type": {
			Description: "Is this a live (green) or idle (blue) stack?",
			Schema: map[string]interface{}{
				"type":     "string",
				"required": true,
				"enum":     []string{"live", "idle"},
			},
		},
		"username": {
			Description: "What username do you use to login to this stack?",
			Schema: map[string]interface{}{
				"type":      "string",
				"required":  true,
				"minLength": 1,
			},
		},
	})
	c.SetDefault("answers", []map[string]interface{}{})
	c.SetDefault("config_templates.default.pem_filepath", "{{.environment}}/{{.username}}-{{.environment}}.pem")
	c.SetDefault("config_templates.default.filepath", `{{.environment}}-{{.stack_name}}-{{.shard_type}}-{{.shard}}{{if ne .username "aws"}}-{{.username}}{{end}}`)
	c.SetDefault("config_templates.default.content", utils.StripIndent(
		`
		ForwardAgent yes
		ForwardX11 no
		HashKnownHosts yes
		IdentitiesOnly yes
		StrictHostKeyChecking no


		Host bastion
		  	Hostname bastion1.{{.shard_type}}.{{.shard}}.{{.stack_name}}{{if ne .environment "prod"}}{{.environment}}{{end}}example.com
		  	User {{if eq .username "aws"}}cloud-user{{else}}{{.username}}{{end}}
		  	IdentityFile {{.template.pem_filepath}}
		  	LocalForward localhost:{{uniquePort .template.filepath}} localhost:8080
		  	UserKnownHostsFile=/dev/null
		  	StrictHostKeyChecking=no

		Host bastion+*
		  	ProxyCommand ssh -F {{.template.filepath}} -W $(echo %h |cut -d+ -f2):%p bastion
		  	User {{if eq .username "aws"}}cloud-user{{else}}{{.username}}{{end}}
		  	IdentityFile {{.template.pem_filepath}}
		  	LogLevel INFO
		  	UserKnownHostsFile=/dev/null
		  	StrictHostKeyChecking=no
	`))
	c.SetDefault("custom_templates", map[string]interface{}{})

	c.SetDefault("pac_template.filepath", `~/drawbridge.pac`)
	c.SetDefault("pac_template.content", utils.StripIndent(
		`
		// This file was automatically generated by Drawbridge
		// Do not modify.

		// Proxy Auto-Config File.
		//

		function FindProxyForURL(url, host){
			//determine if we need to use proxy. 

			{{range .}}
			if(dnsDomainIs(host, ".internal.{{.shard_type}}.{{.shard}}.{{.stack_name}}{{if ne .environment "prod"}}{{.environment}}{{end}}example.com")){
				return "PROXY localhost:{{uniquePort .config.filepath}}";
			}
			{{end}}

			// use default connection (skip proxy)
			return "DIRECT";
		}
	`))

	//if you want to load a non-standard location system config file (~/drawbridge.yml), use ReadConfig
	c.SetConfigType("yaml")
	//c.SetConfigName("drawbridge")
	//c.AddConfigPath("$HOME/")

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
	//}
	//
	//

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
	// TODO: ensure that all custom_templates.filepaths are absolute or start with ~/
	// language=json
	configFileSchema := `
	{
		"type": "object",
		"additionalProperties":false,
		"required": ["version"],
		"properties":{
			"version": {
				"type": "integer"
			},
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
					"active_custom_templates": {
						"type":"array",
						"uniqueItems": true,
						"items":[{"type":"string"}]
					},
					"ui_group_priority": {
						"type":"array",
						"uniqueItems": true,
						"items":[{"type":"string"}],
						"maxItems": 4
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
					"^[a-z0-9\\_]*$":{
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
									"required": {
										"type": "boolean"
									},
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
				"additionalProperties":false,
				"items": {
					"oneOf": [
						{
							"type" : "object",
							"additionalProperties":false,
							"required": ["_file"],
							"properties" : {
                    			"_file" : {
                        			"type" : "string",
									"pattern": "^(/[^/]+)+$"
                    			}
                			}
						},
						{
							"type" : "object",
							"additionalProperties":false,
							"patternProperties": {
								"^[a-z0-9\\_]*$": {
								}
							}
						}
					]
				}
			},
			"variables":{
				"type": "object",
				"patternProperties": {
					"^[a-z0-9]*$":{
						"type":"string"
					}
				}
			},
			"config_templates":{
				"type": "object",
				"patternProperties": {
					"^[a-z0-9]*$":{
						"type":"object",
						"additionalProperties":false,
						"required": ["filepath", "content", "pem_filepath"],
						"properties": {
							"filepath": {
								"type": "string"
							},
							"content": {
								"type": "string"
							},
							"pem_filepath": {
								"type": "string"
							}
						}
					}
				}
			},
			"custom_templates":{
				"type": "object",
				"patternProperties": {
					"^[a-z0-9]*$":{
						"type":"object",
						"additionalProperties":false,
						"required": ["filepath","content"],
						"properties": {
							"filepath": {
								"type": "string"
							},
							"content": {
								"type": "string"
							}
						}
					}
				}
			},
			"pac_template":{
				"type":"object",
				"additionalProperties":false,
				"required": ["filepath","content"],
				"properties": {
					"filepath": {
						"type": "string"
					},
					"content": {
						"type": "string"
					}
				}
			}
		}
	}
	`

	schemaLoader := gojsonschema.NewStringLoader(configFileSchema)
	documentLoader := gojsonschema.NewGoLoader(configContent)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		errorMsg := ""
		for _, err := range result.Errors() {
			// Err implements the ResultError interface
			errorMsg += fmt.Sprintf("- %s\n", err)
		}

		return errors.ConfigValidationError(fmt.Sprintf("There was an error validating this config:\n %v ", errorMsg))
	}
	return nil
}

func (c *configuration) InternalQuestionKeys() []string {
	//list of internal keys, can be filtered out when printing, etc.
	return []string{"config_dir", "pem_dir", "active_config_template", "active_custom_templates", "ui_group_priority", "ui_question_hidden", "custom", "config", "template"}
}

func (c *configuration) GetProvidedAnswerList() ([]map[string]interface{}, error) {
	//deserialize
	answerList := []map[string]interface{}{}
	err := c.UnmarshalKey("answers", &answerList)
	return answerList, err
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

func (c *configuration) GetPacTemplate() (template.PacTemplate, error) {
	//deserialize Template

	template := template.PacTemplate{}
	err := c.UnmarshalKey("pac_template", &template)
	return template, err
}

func (c *configuration) GetConfigTemplates() (map[string]template.ConfigTemplate, error) {
	//deserialize Templates
	templateMap := map[string]template.ConfigTemplate{}
	err := c.UnmarshalKey("config_templates", &templateMap)
	return templateMap, err
}

func (c *configuration) GetActiveConfigTemplate() (template.ConfigTemplate, error) {
	//deserialize Templates
	activeTemplateName := c.GetString("options.active_config_template")

	allTemplates, err := c.GetConfigTemplates()
	if err != nil {
		return template.ConfigTemplate{}, err
	}

	activeTemplate := allTemplates[activeTemplateName]
	return activeTemplate, nil
}

func (c *configuration) GetCustomTemplates() (map[string]template.FileTemplate, error) {
	//deserialize Templates
	templateMap := map[string]template.FileTemplate{}
	err := c.UnmarshalKey("custom_templates", &templateMap)
	return templateMap, err
}

func (c *configuration) GetActiveCustomTemplates() ([]template.FileTemplate, error) {
	//deserialize Templates
	activeTemplateNames := c.GetStringSlice("options.active_custom_templates")

	allTemplates, err := c.GetCustomTemplates()
	if err != nil {
		return nil, err
	}
	activeTemplates := []template.FileTemplate{}

	for _, activeTemplateName := range activeTemplateNames {
		activeTemplate := allTemplates[activeTemplateName]
		activeTemplates = append(activeTemplates, activeTemplate)
	}
	return activeTemplates, nil
}

func (c *configuration) SetOptionsFromAnswers(answerValues map[string]interface{}) {

	// get current options
	options := map[string]interface{}{}
	c.UnmarshalKey("options", &options)

	optionKeys := []string{}
	for key := range options {
		optionKeys = append(optionKeys, key)
	}

	//find optionKeys in answerValues
	for _, optionKey := range optionKeys {
		//check if the key is set as an answer/default
		if answerOptionValue, ok := answerValues[optionKey]; ok {
			//this answer is actualy for an option. lets set it.
			//logger.Debugf("\nSetting option from Answer: %v  (%v)", optionKey, answerOptionValue)
			options[optionKey] = answerOptionValue
		}
	}

	//set the updated options in the config.
	c.Set("options", options)
}
