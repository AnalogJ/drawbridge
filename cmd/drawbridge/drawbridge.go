package main

import (
	"fmt"
	"os"
	"time"

	"drawbridge/pkg/utils"
	"drawbridge/pkg/config"
	"drawbridge/pkg/version"
	"gopkg.in/urfave/cli.v2"
	"drawbridge/pkg/create"
	"drawbridge/pkg/list"
	"drawbridge/pkg/connect"
	"bufio"
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

	flags, err := createFlags(config)
	if err != nil {
		fmt.Printf("FATAL: %+v\n", err)
		os.Exit(1)
	}


	app := &cli.App{
		Name:     "drawbridge",
		Usage:    "Bastion tunneling made easy",
		Version:  version.VERSION,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Jason Kulatunga",
				Email: "jason@thesparktree.com",
			},
		},
		Before: func(c *cli.Context) error {

			drawbridge := "https://blog.thesparktree.com"

			versionInfo := fmt.Sprintf("%s.%s-%s", goos, goarch, version.VERSION)

			subtitle := drawbridge + utils.LeftPad2Len(versionInfo, " ", 53-len(drawbridge))

			fmt.Fprintf(c.App.Writer, fmt.Sprintf(utils.StripIndent(
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
				Action: func(c *cli.Context) error {
					//TODO: list the existing configurations from answers, ask user to specify an exiting config, or create a new one.

					//TODO: check if the user decides to create one from scratch.

					//pass in CLI answer data.
					cliAnswers, err := createFlagAnswers(config, c.FlagNames(), c)
					if err != nil {
						return err
					}

					createEngine := create.CreateEngine{ Config: config }
					return createEngine.Start(cliAnswers)
				},

				Flags: flags,
			},
			{
				Name:  "list",
				Usage: "List all drawbridge managed ssh configs",
				Action: func(c *cli.Context) error {


					listEngine := list.ListEngine{ Config: config }
					return listEngine.Start()
				},

				Flags: flags,
			},
			{
				Name:  "connect",
				Usage: "Connect to a drawbridge managed ssh config",
				Action: func(c *cli.Context) error {

					listEngine := list.ListEngine{ Config: config }
					listEngine.Start()

					reader := bufio.NewReader(os.Stdin)
					text, _ := reader.ReadString('\n')
					fmt.Println(text)
					text = strings.TrimSpace(text)
					i, err := strconv.Atoi(text)
					if err != nil{
						return err
					}
					answerIndex := i - 1

					connectEngine := connect.ConnectEngine{ Config: config }
					return connectEngine.Start(listEngine.OrderedAnswers[answerIndex].(map[string]interface{}))
				},

				Flags: flags,
			},
		},
	}

	app.Run(os.Args)
}

func createFlags(appConfig config.Interface) ([]cli.Flag, error){
	flags := []cli.Flag{
		&cli.StringFlag{
			Name: "active_config_template",
			Usage: "Active config_template",
			Value: appConfig.GetString("options.active_config_template"),
		},
		&cli.StringSliceFlag{
			Name: "active_extra_templates",
			Usage: "Activated extra_templates",
			Value: cli.NewStringSlice(appConfig.GetStringSlice("options.active_extra_templates")...),
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
				Name:    k,
				Usage:   v.Description,
			}
			defaultValue, ok := v.DefaultValue.(string)
			if ok {
				newFlag.Value = defaultValue
			}

			flags = append(flags, newFlag)
		} else if questionType == "integer" {
			newFlag := &cli.IntFlag{
				Name:    k,
				Usage:   v.Description,
			}
			defaultValue, ok := v.DefaultValue.(int)
			if ok {
				newFlag.Value = defaultValue
			}

			flags = append(flags, newFlag)
		} else if questionType == "boolean" {
			newFlag := &cli.BoolFlag{
				Name:    k,
				Usage:   v.Description,
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

func createFlagAnswers(appConfig config.Interface, cliFlags []string, c *cli.Context) (map[string]interface{}, error){

	cliAnswers := map[string]interface{}{}

	for _, questionKey := range cliFlags {
		question, err := appConfig.GetQuestion(questionKey)
		if err != nil{
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
