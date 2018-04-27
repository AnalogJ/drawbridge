package project

type projectData struct {

	//required for all Projects
	Answers map[string]interface{}

	// Only populated for Projects that have actually been created (stored to disk)
	AnswerFilePath          string
	ConfigFilePath          string
	PemFilePath             string
	CustomTemplateFilePaths []string
}
