package utils

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"regexp"
	"strings"
	"syscall"
)

func StdinQueryPassword(question string) (string, error) {

	fmt.Println(color.BlueString(question))
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	text := strings.TrimSpace(string(bytePassword))
	return text, nil
}

func StdinQuery(question string) string {

	fmt.Println(color.BlueString(question))
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}

func StdinQueryBoolean(question string) bool {

	text := StdinQuery(question)
	text = strings.ToLower(text)

	if text == "true" || text == "yes" || text == "y" {
		return true
	} else if text == "false" || text == "no" || text == "n" {
		return false
	} else {
		color.Yellow("WARNING: invalid response only true/yes/y/false/no/n allowed not `%v`.\nAssuming `no`", text)
		return false
	}
}

func StdinQueryInt(question string) (int, error) {

	text := StdinQuery(question)
	return StringToInt(text)
}

func StdinQueryRegex(message string, regexPattern string, friendlyPattern string) string {
	for true {

		//prompt the user to enter a valid choice
		text := StdinQuery(fmt.Sprintf("%v (%s):", message, friendlyPattern))

		isValid, err := regexp.MatchString(regexPattern, text)

		if err != nil {
			color.HiRed("ERROR: %v", err)
			continue
		} else if !isValid {
			color.HiRed("Invalid alias. Must match pattern: %s", friendlyPattern)
			continue
		}

		return text
	}
	return ""
}
