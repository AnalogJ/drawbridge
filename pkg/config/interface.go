package config

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
	UnmarshalKey(key string, rawVal interface{}) error

	GetQuestion(questionKey string) (Question, error)
	GetQuestions() (map[string]Question, error)
	//GetQuestionsSchema() (map[string]interface{}, error)
	//GetQuestionSchema(question Question) (map[string]interface{}, error)

	GetConfigTemplates() (map[string]Template, error)
	GetActiveConfigTemplate() (Template, error)
	GetExtraTemplates() (map[string]Template, error)
	GetActiveExtraTemplates() ([]Template, error)
}
