package config_test

import (
	"drawbridge/pkg/config"
	"github.com/stretchr/testify/require"
	"path"
	"testing"
)

func TestConfiguration_init_ShouldCorrectlyInitializeConfiguration(t *testing.T) {
	t.Parallel()

	//test
	testConfig, err := config.Create()

	//assert
	require.NoError(t, err, "should not have an error")
	require.Equal(t, "~/.ssh/drawbridge", testConfig.GetString("options.config_dir"), "should populate config_dir with default")
	require.Equal(t, "~/.ssh", testConfig.GetString("options.pem_dir"), "should populate pem_dir with default")
	require.Equal(t, "default", testConfig.GetString("options.active_config_template"), "should populate active_config_template with default")
	require.Equal(t, []string{}, testConfig.GetStringSlice("options.active_extra_templates"), "should populate active_config_template with empty list")
}

func TestConfiguration_ReadConfig_InvalidFilePath(t *testing.T) {
	t.Parallel()
	//setup
	testConfig, _ := config.Create()
	err := testConfig.ReadConfig(path.Join("does", "not", "exist.yml"))

	//assert
	require.Error(t, err, "should raise an error")
}

func TestConfiguration_ReadConfig_Simple(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, err := config.Create()
	require.NoError(t, err, "should create valid default config file")

	//test
	err = testConfig.ReadConfig(path.Join("testdata", "valid_simple_config.yaml"))

	//assert
	require.NoError(t, err, "should be valid config file")
	require.Equal(t, "~/.ssh/drawbridge", testConfig.GetString("options.config_dir"), "should populate config_dir with default")
	require.Equal(t, []string{"default", "knife"}, testConfig.GetStringSlice("options.active_extra_templates"), "should populate active_extra_templates with overrides")
	require.Equal(t, "~/.ssh/drawbridge/pem", testConfig.GetString("options.pem_dir"), "should populate pem_dir with overrides")

}

//func TestConfiguration_ReadConfig_Answers(t *testing.T) {
//	t.Parallel()
//
//	//setup
//	testConfig, _ := config.Create()
//
//	//test
//	testConfig.ReadConfig(path.Join("testdata", "simple_config.yaml"))
//
//
//	//assert
//	require.Equal(t, "~/.ssh/drawbridge", testConfig.GetString("options.config_dir"), "should populate config_dir with default")
//	require.Equal(t, "~/.ssh/drawbridge/pem", testConfig.GetString("options.pem_dir"), "should populate pem_dir with default")
//}
//
//func TestConfiguration_ReadConfig_AnswersFile(t *testing.T) {
//	t.Parallel()
//
//	//setup
//	testConfig, _ := config.Create()
//
//	//test
//	testConfig.ReadConfig(path.Join("testdata", "simple_config.yaml"))
//
//
//	//assert
//	require.Equal(t, "~/.ssh/drawbridge", testConfig.GetString("options.config_dir"), "should populate config_dir with default")
//	require.Equal(t, "~/.ssh/drawbridge/pem", testConfig.GetString("options.pem_dir"), "should populate pem_dir with default")
//}

func TestConfiguration_ReadConfig_Questions(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "valid_questions.yaml"))

	//assert
	require.NoError(t, err, "should correctly parse config file.")
}

func TestConfiguration_ReadConfig_QuestionsWithMissingTypeReturnsError(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "invalid_questions_missing_type.yaml"))

	//assert
	require.Error(t, err, "should return an error if the question type is missing.")
}

func TestConfiguration_ReadConfig_QuestionsWithEmptyTypeReturnsError(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "invalid_questions_empty_type.yaml"))

	//assert
	require.Error(t, err, "should return an error if the question type is empty.")
}

func TestConfiguration_ReadConfig_UnsupportedTopLevelKey(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "invalid_unsupported_top_level_key.yaml"))

	//assert
	require.Error(t, err, "should return an error if there is an unsupported top level key.")
}

func TestConfiguration_ReadConfig_DuplicateActiveTemplates(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "invalid_active_templates.yaml"))

	//assert
	require.Error(t, err, "should return an error if there is an duplicate active template")
}

func TestConfiguration_ReadConfig_OverrideDefaultConfigTemplate(t *testing.T) {
	t.Parallel()

	//setup
	testConfig, _ := config.Create()

	//test
	err := testConfig.ReadConfig(path.Join("testdata", "valid_config_template.yaml"))
	require.NoError(t, err, "should allow overriding default config template.")

	configTmpl, err := testConfig.GetActiveConfigTemplate()

	//assert
	require.NoError(t, err, "should allow overriding default config template.")
	require.Equal(t,"{{.environment}}-{{.username}}", configTmpl.FilePath)

}