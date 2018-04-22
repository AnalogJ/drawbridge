package template

type Template struct {
	Content string `mapstructure:"content"`

	data map[string]interface{}
}
