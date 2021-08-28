package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"time"

	"github.com/analogj/drawbridge/pkg/actions"
	"github.com/analogj/drawbridge/pkg/config"
	"github.com/analogj/drawbridge/pkg/errors"
	"github.com/analogj/drawbridge/pkg/project"
	"github.com/analogj/drawbridge/pkg/utils"
	"github.com/analogj/drawbridge/pkg/version"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"strings"
)

var goos string
var goarch string

func main() {

	config, err := config.Create()
	if err != nil {
		fmt.Printf("FATAL: %+v\n", err)
		os.Exit(1)
	}

	//we're going to load the config file manually, since we need to validate it.
	err = config.ReadConfig("~/drawbridge.yaml")          // Find and read the config file
	if _, ok := err.(errors.ConfigFileMissingError); ok { // Handle errors reading the config file
		//ignore "could not find config file"
	} else if err != nil {
		os.Exit(1)
	}

	createFlags, err := createFlags(config)
	if err != nil {
		fmt.Printf("FATAL: %+v\n", err)
		os.Exit(1)
	}

	cli.CommandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}
USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}
CATEGORY:
   {{.Category}}{{end}}{{if .Description}}
DESCRIPTION:
   {{.Description}}{{end}}{{if .VisibleFlags}}
OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

	app := &cli.App{
		Name:     "drawbridge",
		Usage:    "Bastion/Jumphost tunneling made easy",
		Version:  version.VERSION,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Jason Kulatunga",
				Email: "jason@thesparktree.com",
			},
		},
		Before: func(c *cli.Context) error {

			drawbridge := "github.com/AnalogJ/drawbridge"

			var versionInfo string
			if len(goos) > 0 && len(goarch) > 0 {
				versionInfo = fmt.Sprintf("%s.%s-%s", goos, goarch, version.VERSION)
			} else {
				versionInfo = fmt.Sprintf("dev-%s", version.VERSION)
			}

			subtitle := drawbridge + utils.LeftPad2Len(versionInfo, " ", 65-len(drawbridge))

			color.New(color.FgGreen).Fprintf(c.App.Writer, fmt.Sprintf(utils.StripIndent(
				`
			 ____  ____    __    _    _  ____  ____  ____  ____    ___  ____
			(  _ \(  _ \  /__\  ( \/\/ )(  _ \(  _ \(_  _)(  _ \  / __)( ___)
			 )(_) ))   / /(__)\  )    (  ) _ < )   / _)(_  )(_) )( (_-. )__)
			(____/(_)\_)(__)(__)(__/\__)(____/(_)\_)(____)(____/  \___/(____)
			%s

			`), subtitle))

			return nil
		},

		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a drawbridge managed ssh config & associated files",
				//UsageText:   "doo - does the dooing",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)
					if c.Bool("debug") {
						log.SetLevel(log.DebugLevel)
					} else {
						log.SetLevel(log.InfoLevel)
					}

					projectList, err := project.CreateProjectListFromProvidedAnswers(config)
					if err != nil {
						return err
					}

					answerData := map[string]interface{}{}
					if projectList.Length() > 0 && utils.StdinQueryBoolean(fmt.Sprintf("Would you like to create a Drawbridge config using preconfigured answers? (%v available). [yes/no]", projectList.Length())) {

						answerData, _, err = projectList.Prompt("Enter number to base your configuration from")
						if err != nil {
							return err
						}
					}

					//extend current answerData with CLI provided options.
					cliAnswers, err := createFlagHandler(config, answerData, c.FlagNames(), c)
					if err != nil {
						return err
					}

					createAction := actions.CreateAction{Config: config}
					return createAction.Start(cliAnswers, c.Bool("dryrun"))
				},

				Flags: createFlags,
			},
			{
				Name:      "list",
				Usage:     "List all drawbridge managed ssh configs",
				ArgsUsage: "[config_number/alias]",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					projectList, err := project.CreateProjectListFromConfigDir(config)
					if err != nil {
						return err
					}

					var answerData map[string]interface{}
					if c.NArg() > 0 {
						answerData, _, err = projectList.GetWithAliasOrIndex(c.Args().Get(0))
						if err != nil {
							return err
						}

					} else {
						answerData, _, err = projectList.Prompt("Enter drawbridge config number or alias to retrieve full info")
						if err != nil {
							return err
						}
					}

					fmt.Print("\nAnswer Data:\n")
					for k, v := range answerData {
						fmt.Printf("\t%v: %v\n", color.YellowString(k), v)
					}

					return nil
				},
				Flags: nil,
			},
			{
				Name:      "connect",
				Usage:     "Connect to a drawbridge managed ssh config",
				ArgsUsage: "[config_number/alias] [dest_server_hostname]",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					projectList, err := project.CreateProjectListFromConfigDir(config)
					if err != nil {
						return err
					}

					var answerData map[string]interface{}
					if c.NArg() > 0 {

						answerData, _, err = projectList.GetWithAliasOrIndex(c.Args().Get(0))
						if err != nil {
							return err
						}

					} else {
						answerData, _, err = projectList.Prompt("Enter drawbridge config number to connect to")
						if err != nil {
							return err
						}
					}

					var destServer string
					if c.IsSet("dest") {
						destServer = c.String("dest")
					} else if c.NArg() >= 2 {
						destServer = c.Args().Get(1)
					} else {
						destServer = ""
					}

					config.SetOptionsFromAnswers(answerData)
					connectAction := actions.ConnectAction{Config: config}
					return connectAction.Start(answerData, destServer, c.Bool("debug"))
				},

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "dest",
						Usage: "Specify the `hostname` of the destination/internal server you would like to connect to.",
					},
					&cli.BoolFlag{
						Name:  "debug",
						Value: false,
						Usage: "Debug mode",
					},
				},
			},
			{
				Name:      "alias",
				Usage:     "Create a named alias for a drawbridge config",
				ArgsUsage: "[config_number] [alias]",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					projectList, err := project.CreateProjectListFromConfigDir(config)
					if err != nil {
						return err
					}

					var answerData map[string]interface{}
					var answerIndex int
					if c.NArg() > 0 {

						answerData, _, err = projectList.GetWithAliasOrIndex(c.Args().Get(0))
						if err != nil {
							return err
						}

					} else {
						answerData, answerIndex, err = projectList.Prompt("Enter drawbridge config number to create alias for")
						if err != nil {
							return err
						}
					}

					fmt.Print("\nAnswer Data:\n")
					for k, v := range answerData {
						fmt.Printf("\t%v: %v\n", color.YellowString(k), v)
					}

					//get the alias name (if provided)
					var configAlias string
					if c.NArg() >= 2 {
						configAlias = c.Args().Get(1)
						isValid, err := regexp.MatchString(`^[A-Za-z][\w-\.]+$`, configAlias)
						if err != nil || !isValid {
							configAlias = utils.StdinQueryRegex("Please provide an alias for the configuration above", `^[A-Za-z][\w-\.]+$`, "a-zA-Z0-9-_.")
						}
					} else {
						configAlias = utils.StdinQueryRegex("Please provide an alias for the configuration above", `^[A-Za-z][\w-\.]+$`, "a-zA-Z0-9-_.")
					}

					_, _, aliasErr := projectList.GetWithAlias(configAlias)
					if aliasErr == nil {
						return errors.ConfigValidationError("alias already exists.")
					}

					color.HiBlue("Setting alias (%s) for config (%d)\n", configAlias, answerIndex+1)

					_, err = projectList.SetAliasForIndex(answerIndex, configAlias)

					return err
				},
			},
			{
				Name:      "download",
				Aliases:   []string{"scp"},
				Usage:     "Download a file from an internal server using drawbridge managed ssh config, syntax is similar to scp command. ",
				ArgsUsage: "[config_number/alias] destination_hostname:remote_filepath local_filepath",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					projectList, err := project.CreateProjectListFromConfigDir(config)
					if err != nil {
						return err
					}

					// PARSE ARGS
					if c.NArg() < 2 || c.NArg() > 3 {
						return errors.InvalidArgumentsError(fmt.Sprintf("2 or 3 arguments required. %v provided", c.Args().Len()))
					}

					index := 0
					strRemoteHostname := ""
					strRemotePath := ""
					strLocalPath := ""

					args := c.Args().Slice()

					var answerData map[string]interface{}

					if c.NArg() == 3 {
						answerData, index, err = projectList.GetWithAliasOrIndex(c.Args().Get(0))
						if err != nil {
							return err
						}
						args = c.Args().Tail()
					}

					remoteParts := strings.Split(args[0], ":")
					if len(remoteParts) != 2 {
						return errors.InvalidArgumentsError(fmt.Sprintf("Invalid `destination_hostname:remote path` format: %s", remoteParts))
					} else {
						strRemoteHostname = remoteParts[0]
						strRemotePath = remoteParts[1]
					}

					strLocalPath = args[1]

					//index is unset, lets prompt for the answerData
					if index == 0 {
						answerData, _, err = projectList.Prompt("Enter number of drawbridge config you would like to download from")
						if err != nil {
							return err
						}
					}

					config.SetOptionsFromAnswers(answerData)
					downloadAction := actions.DownloadAction{Config: config}
					return downloadAction.Start(answerData, strRemoteHostname, strRemotePath, strLocalPath)
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete drawbridge managed ssh config(s)",
				ArgsUsage: "[config_number]",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					projectList, err := project.CreateProjectListFromConfigDir(config)
					if err != nil {
						return err
					}

					var answerData map[string]interface{}

					if c.Bool("all") {
						//check if the user wants to delete all configs
						deleteAction := actions.DeleteAction{Config: config}
						return deleteAction.All(projectList.GetAll(), c.Bool("force"))

					} else if c.NArg() > 0 {
						//check if the user specified a config number in the args.

						answerData, _, err = projectList.GetWithAliasOrIndex(c.Args().Get(0))
						if err != nil {
							return err
						}

					} else {
						// prompt the user to determine which configs to delete.
						answerData, _, err = projectList.Prompt("Enter drawbridge config number to delete")
						if err != nil {
							return err
						}
					}

					//delete one config file.
					config.SetOptionsFromAnswers(answerData)
					deleteAction := actions.DeleteAction{Config: config}
					err = deleteAction.One(answerData, c.Bool("force"))

					if err != nil {
						//print an error message here:
						return err
					} else {
						color.Green("Finished")
						return nil
					}
				},

				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "force",
						Usage: "Force delete with no confirmation",
					},
					&cli.BoolFlag{
						Name:  "all",
						Usage: "Delete all configuration files. ",
					},

					//TODO: add dry run support
				},
			},
			{
				Name:  "proxy",
				Usage: "Build/Rebuild a Proxy auto-config (PAC) file to access websites through Drawbridge tunnels",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					projectList, err := project.CreateProjectListFromConfigDir(config)
					if err != nil {
						return err
					}
					answerDataList := projectList.GetAll()

					proxyAction := actions.ProxyAction{Config: config}
					return proxyAction.Start(answerDataList, false)
				},
			},
			{
				Name:  "update",
				Usage: "Update drawbridge to the latest version",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					if len(goos) == 0 && len(goarch) == 0 {
						//dev mode,
						color.Yellow("WARNING: Binary was built from source, not released. Auto-update may not work correctly")
					}

					updateAction := actions.UpdateAction{Config: config}
					return updateAction.Start()
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(color.HiRedString("ERROR: %v", err))
	}

}

func createFlags(appConfig config.Interface) ([]cli.Flag, error) {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "active_config_template",
			Usage:       "Active config_template",
			DefaultText: appConfig.GetString("options.active_config_template"),
		},
		&cli.StringSliceFlag{
			Name:        "active_custom_templates",
			Usage:       "Activated custom_templates",
			DefaultText: strings.Join(appConfig.GetStringSlice("options.active_custom_templates"), ", "),
		},
		&cli.BoolFlag{
			Name:  "dryrun",
			Usage: "Dry Run mode. Will print files and paths to STDOUT rather than writing them to disk.",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable verbose logging for debugging",
			Value: false,
		},
	}

	configQuestions, err := appConfig.GetQuestions()
	if err != nil {
		return nil, err
	}
	for k, v := range configQuestions {
		questionType := v.GetType()

		if questionType == "string" {
			newFlag := &cli.StringFlag{
				Name:  k,
				Usage: v.Description,
			}
			defaultValue, ok := v.DefaultValue.(string)
			if ok {
				newFlag.Value = defaultValue
			}

			flags = append(flags, newFlag)
		} else if questionType == "integer" {
			newFlag := &cli.IntFlag{
				Name:  k,
				Usage: v.Description,
			}
			defaultValue, ok := v.DefaultValue.(int)
			if ok {
				newFlag.Value = defaultValue
			}

			flags = append(flags, newFlag)
		} else if questionType == "boolean" {
			newFlag := &cli.BoolFlag{
				Name:  k,
				Usage: v.Description,
			}
			defaultValue, ok := v.DefaultValue.(bool)
			if ok {
				newFlag.Value = defaultValue
			}

			flags = append(flags, newFlag)
		}
	}
	return flags, nil
}

func createFlagHandler(appConfig config.Interface, answerValues map[string]interface{}, cliFlags []string, c *cli.Context) (map[string]interface{}, error) {
	//there's 4 special cases we need to handle for "defaultOptions":
	//case 1: no flag override and no answer option
	//case 2: flag override and no answer option
	//case 3: no flag override and answer option
	//case 4: flag override and answer option

	// the override order is:
	// default value from config
	// answer value from config
	// flag override

	//get default defaultOptions from the config
	options := map[string]interface{}{}
	appConfig.UnmarshalKey("options", &options)
	log.Debugf("\nDefault Options: %v", options)

	optionKeys := []string{}
	for key := range options {
		optionKeys = append(optionKeys, key)
	}

	cliAnswers := answerValues

	//find optionKeys in answerValues
	for _, optionKey := range optionKeys {
		//check if the key is set as an answer/default
		if answerOptionValue, ok := answerValues[optionKey]; ok {
			//this answer is actualy for an option. lets set it.
			log.Debugf("\nSetting option from Answer: %v  (%v)", optionKey, answerOptionValue)
			options[optionKey] = answerOptionValue
			//appConfig.SetDefault(fmt.Sprintf("options.%v", optionKey), answerOptionValue)
		}
	}

	for _, flagName := range cliFlags {

		if utils.SliceIncludes(optionKeys, flagName) {
			//this flag is actually an "option". Lets set it.
			log.Debugf("\nSetting option from CLI: %v (%v)", flagName, c.String(flagName))
			options[flagName] = c.String(flagName)
			//appConfig.SetDefault(fmt.Sprintf("options.%v", flagName), c.String(flagName))
			continue
		}

		//skip dryrun & debug
		if flagName == "dryrun" || flagName == "debug" {
			continue
		}

		questionKey := flagName

		question, err := appConfig.GetQuestion(questionKey)
		if err != nil {
			return nil, err
		}

		questionType := question.GetType()

		if questionType == "string" {
			cliAnswers[questionKey] = c.String(questionKey)

		} else if questionType == "integer" {
			cliAnswers[questionKey] = c.Int(questionKey)

		} else if questionType == "boolean" {
			cliAnswers[questionKey] = c.Bool(questionKey)
		}
	}

	//set the config options section
	appConfig.Set("options", options)

	if log.GetLevel() == log.DebugLevel {
		afterOptions := map[string]interface{}{}
		appConfig.UnmarshalKey("options", &afterOptions)
		log.Debugf("\nOptions after overrides: %v", afterOptions)
	}

	return cliAnswers, nil
}
