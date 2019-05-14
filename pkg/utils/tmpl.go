package utils

import (
	"bytes"
	"encoding/json"
	"hash/fnv"
	"strings"
	"text/template"
)

func PopulatePathTemplate(pathTmplContent string, data interface{}) (string, error) {
	tmplFilepath, err := PopulateTemplate(pathTmplContent, data)
	if err != nil {
		return "", nil
	}
	tmplFilepath, err = ExpandPath(tmplFilepath)
	if err != nil {
		return "", nil
	}
	return tmplFilepath, nil
}

func PopulateTemplate(tmplContent string, data interface{}) (string, error) {
	//set functions
	fns := template.FuncMap{
		"uniquePort":          UniquePort,
		"expandPath":          ExpandPath,
		"stringsCompare":      strings.Compare,
		"stringsContains":     strings.Contains,
		"stringsContainsAny":  strings.ContainsAny,
		"stringsCount":        strings.Count,
		"stringsEqualFold":    strings.EqualFold,
		"stringsHasPrefix":    strings.HasPrefix,
		"stringsHasSuffix":    strings.HasSuffix,
		"stringsIndex":        strings.Index,
		"stringsIndexAny":     strings.IndexAny,
		"stringsJoin":         strings.Join,
		"stringsLastIndex":    strings.LastIndex,
		"stringsLastIndexAny": strings.LastIndexAny,
		"stringsRepeat":       strings.Repeat,
		"stringsReplace":      strings.Replace,
		"stringsSplit":        strings.Split,
		"stringsSplitAfter":   strings.SplitAfter,
		"stringsSplitAfterN":  strings.SplitAfterN,
		"stringsSplitN":       strings.SplitN,
		"stringsTitle":        strings.Title,
		"stringsToLower":      strings.ToLower,
		"stringsToTitle":      strings.ToTitle,
		"stringsToUpper":      strings.ToUpper,
		"stringsTrim":         strings.Trim,
		"stringsTrimLeft":     strings.TrimLeft,
		"stringsTrimPrefix":   strings.TrimPrefix,
		"stringsTrimRight":    strings.TrimRight,
		"stringsTrimSpace":    strings.TrimSpace,
		"stringsTrimSuffix":   strings.TrimSuffix,
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
func UniquePort(data interface{}) (int, error) {

	var contentData []byte
	switch in := data.(type) {
	case string:
		contentData = []byte(in)
	default:
		jsonData, err := json.Marshal(StringifyYAMLMapKeys(in))
		if err != nil {
			return 0, err
		}
		contentData = jsonData
	}

	hash := fnv.New32a()
	hash.Write(contentData)

	//last port - last privileged port.
	portRange := 65535 - 1023

	uniquePort := (hash.Sum32() % uint32(portRange)) + 1023
	return int(uniquePort), nil
}
