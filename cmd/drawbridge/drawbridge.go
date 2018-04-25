package main

import (
	"fmt"
	"os"
	"time"

	"drawbridge/pkg/actions"
	"drawbridge/pkg/config"
	"drawbridge/pkg/errors"
	"drawbridge/pkg/utils"
	"drawbridge/pkg/version"
	"github.com/fatih/color"
	"gopkg.in/urfave/cli.v2"
	"log"
	"strconv"
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
	err = config.ReadConfig("~/drawbridge.yaml")              // Find and read the config file
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

	app := &cli.App{
		Name:     "drawbridge",
		Usage:    "Bastion/Jumphost tunneling made easy",
		Version:  version.VERSION,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Jason Kulatunga",
				Email: "jason@thesparktree.com",
			},
		},
		Before: func(c *cli.Context) error {

			drawbridge := "github.com/AnalogJ/drawbridge"

			versionInfo := fmt.Sprintf("%s.%s-%s", goos, goarch, version.VERSION)

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
					//TODO: list the existing configurations from answers, ask user to specify an exiting config, or create a new one.

					//TODO: check if the user decides to create one from scratch.

					//pass in CLI answer data.
					cliAnswers, err := createFlagHandler(config, c.FlagNames(), c)
					if err != nil {
						return err
					}

					createAction := actions.CreateAction{Config: config}
					return createAction.Start(cliAnswers, c.Bool("dryrun"))
				},

				Flags: createFlags,
			},
			{
				Name:  "list",
				Usage: "List all drawbridge managed ssh configs",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					listAction := actions.ListAction{Config: config}
					return listAction.Start()
				},
			},
			{
				Name:  "connect",
				Usage: "Connect to a drawbridge managed ssh config",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					listAction := actions.ListAction{Config: config}
					listAction.Start()

					if len(listAction.OrderedAnswers) == 0 {
						return nil
					}

					var answerIndex int

					if c.IsSet("drawbridge_id") {
						answerIndex = c.Int("drawbridge_id")
					} else {
						text := utils.StdinQuery(fmt.Sprintf("Enter number of drawbridge config you would like to connect to (%v-%v):", 1, len(listAction.OrderedAnswers)))
						i, err := strconv.Atoi(text)
						if err != nil {
							return err
						}
						answerIndex = i - 1
						maxAnswerIndex := len(listAction.OrderedAnswers)
						if answerIndex >= maxAnswerIndex {
							return errors.AnswerValidationError(fmt.Sprintf("Invalid selection. Please enter a number from 1-%v", maxAnswerIndex))
						}
					}

					connectAction := actions.ConnectAction{Config: config}
					return connectAction.Start(listAction.OrderedAnswers[answerIndex].(map[string]interface{}), c.String("dest"))
				},

				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "drawbridge_id",
						Usage: "Specify the drawbridge configuration to use",
					},
					&cli.StringFlag{
						Name:  "dest",
						Usage: "Specify the `hostname` of the destination/internal server you would like to connect to.",
					},
				},
			},
			{
				Name:      "download",
				Aliases:   []string{"scp"},
				Usage:     "Download a file from an internal server using drawbridge managed ssh config, syntax is similar to scp command. ",
				ArgsUsage: "destination_hostname:remote_filepath local_filepath",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					if c.NArg() != 2 {
						return cli.Exit(fmt.Sprintf("Invalid, 2 arguments required: %s", c.Args()), 1)
					}

					remoteParts := strings.Split(c.Args().First(), ":")
					if len(remoteParts) != 2 {
						return cli.Exit(fmt.Sprintf("Invalid, please specify destination hostname and remote path: %s", remoteParts), 1)
					}

					listAction := actions.ListAction{Config: config}
					listAction.Start()

					var answerIndex int

					if c.IsSet("drawbridge_id") {
						answerIndex = c.Int("drawbridge_id")
					} else {

						text := utils.StdinQuery(fmt.Sprintf("Enter number of drawbridge config you would like to download from (%v-%v):", 1, len(listAction.OrderedAnswers)))
						i, err := strconv.Atoi(text)
						if err != nil {
							return err
						}
						answerIndex = i - 1
						maxAnswerIndex := len(listAction.OrderedAnswers)
						if answerIndex >= maxAnswerIndex {
							return errors.AnswerValidationError(fmt.Sprintf("Invalid selection. Please enter a number from 1-%v", maxAnswerIndex))
						}
					}

					downloadAction := actions.DownloadAction{Config: config}

					return downloadAction.Start(listAction.OrderedAnswers[answerIndex].(map[string]interface{}), remoteParts[0], remoteParts[1], c.Args().Get(1))
				},

				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "drawbridge_id",
						Usage: "Specify the drawbridge configuration to use",
					},
				},
			},
			{
				Name:  "delete",
				Usage: "Delete drawbridge managed ssh config(s)",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					listAction := actions.ListAction{Config: config}

					if c.Bool("all") {
						deleteAction := actions.DeleteAction{Config: config}
						answersList, err := listAction.RenderedAnswersList()
						if err != nil {
							return nil
						}
						return deleteAction.All(answersList, c.Bool("force"))
					} else {
						listAction := actions.ListAction{Config: config}
						listAction.Start()

						var answerIndex int

						if c.IsSet("drawbridge_id") {
							answerIndex = c.Int("drawbridge_id")
						} else {

							text := utils.StdinQuery(fmt.Sprintf("Enter number of drawbridge config you would like to delete (%v-%v):", 1, len(listAction.OrderedAnswers)))
							i, err := strconv.Atoi(text)
							if err != nil {
								return err
							}
							answerIndex = i - 1

							maxAnswerIndex := len(listAction.OrderedAnswers)
							if answerIndex >= maxAnswerIndex {
								return errors.AnswerValidationError(fmt.Sprintf("Invalid selection. Please enter a number from 1-%v", maxAnswerIndex))
							}
						}

						deleteAction := actions.DeleteAction{Config: config}
						err := deleteAction.One(listAction.OrderedAnswers[answerIndex].(map[string]interface{}), c.Bool("force"))

						if err != nil {
							//print an error message here:
							return err
						} else {
							color.Green("Finished")
							return nil
						}
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
					&cli.IntFlag{
						Name:  "drawbridge_id",
						Usage: "Specify the drawbridge configuration to delete",
					},

					//TODO: add dry run support
				},
			},
			{
				Name:  "proxy",
				Usage: "Build/Rebuild a Proxy auto-config (PAC) file to access websites through Drawbridge tunnels",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					listAction := actions.ListAction{Config: config}
					answerDataList, err := listAction.RenderedAnswersList()
					if err != nil {
						return err
					}

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
			Name:  "active_config_template",
			Usage: "Active config_template",
			Value: appConfig.GetString("options.active_config_template"),
		},
		&cli.StringSliceFlag{
			Name:  "active_custom_templates",
			Usage: "Activated custom_templates",
			Value: cli.NewStringSlice(appConfig.GetStringSlice("options.active_custom_templates")...),
		},
		&cli.BoolFlag{
			Name:  "dryrun",
			Usage: "Dry Run mode. Will print files and paths to STDOUT rather than writing them to disk.",
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

func createFlagHandler(appConfig config.Interface, cliFlags []string, c *cli.Context) (map[string]interface{}, error) {

	cliAnswers := map[string]interface{}{}

	for _, flagName := range cliFlags {
		//handle options
		options := map[string]interface{}{}
		appConfig.UnmarshalKey("options", &options)
		if _, ok := options[flagName]; ok {
			//this flag is actually an option. Lets set it.
			appConfig.Set(fmt.Sprintf("options.%v", flagName), c.String(flagName))
			continue
		}

		//skip dryrun
		if flagName == "dryrun" {
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

	return cliAnswers, nil
}
