package config

import (
	"github.com/analogj/drawbridge/pkg/config/template"
	"github.com/spf13/viper"
)

// Create mock using:
// mockgen -source=pkg/config/interface.go -destination=pkg/config/mock/mock_config.go
type Interface interface {
	Init() error
	ReadConfig(configFilePath string) error
	Set(key string, value interface{})
	SetDefault(key string, value interface{})
	AllSettings() map[string]interface{}
	IsSet(key string) bool
	Get(key string) interface{}
	GetBool(key string) bool
	GetInt(key string) int
	GetString(key string) string
	GetStringSlice(key string) []string
	UnmarshalKey(key string, rawVal interface{}, decoderOpts ...viper.DecoderConfigOption) error

	GetProvidedAnswerList() ([]map[string]interface{}, error)
	InternalQuestionKeys() []string
	GetQuestion(questionKey string) (Question, error)
	GetQuestions() (map[string]Question, error)
	//GetQuestionsSchema() (map[string]interface{}, error)
	//GetQuestionSchema(question Question) (map[string]interface{}, error)

	GetPacTemplate() (template.PacTemplate, error)
	GetConfigTemplates() (map[string]template.ConfigTemplate, error)
	GetActiveConfigTemplate() (template.ConfigTemplate, error)
	GetCustomTemplates() (map[string]template.FileTemplate, error)
	GetActiveCustomTemplates() ([]template.FileTemplate, error)
}
