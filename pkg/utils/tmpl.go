package utils

import (
	"bytes"
	"encoding/json"
	"hash/fnv"
	"text/template"
)

func PopulateTemplate(tmplContent string, data map[string]interface{}) (string, error) {
	//set functions
	fns := template.FuncMap{
		"uniquePort": UniquePort,
		"expandPath": ExpandPath,
	}

	// prep the template, set the option
	tmpl, err := template.New("populate").Option("missingkey=error").Funcs(fns).Parse(tmplContent)
	if err != nil {
		return "", err
	}

	//specify that any missing keys in the template will throw an error
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	//convert buffered content to string
	return buf.String(), nil
}

// https://play.golang.org/p/k8bws03uid
func UniquePort(data map[string]interface{}) (int, error) {
	jsonString, err := json.Marshal(StringifyYAMLMapKeys(data))
	if err != nil {
		return 0, err
	}

	hash := fnv.New32a()
	hash.Write([]byte(jsonString))

	//last port - last privileged port.
	portRange := 65535 - 1023

	uniquePort := (hash.Sum32() % uint32(portRange)) + 1023
	return int(uniquePort), nil
}
