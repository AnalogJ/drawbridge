package actions

import (
	"drawbridge/pkg/config"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"drawbridge/pkg/version"
	"github.com/inconshreveable/go-update"
	"runtime"
	"drawbridge/pkg/errors"
)

type UpdateAction struct {
	Config config.Interface
}

type GithubReleaseInfo struct {
	TagName string `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
	Assets []struct {
		Url string `json:"url"`
		Name string `json:"name"`
	} `json:"assets"`
}


// https://github.com/hyperhq/hypercli/blob/302a6b530148f6a777cd6b8772f706ab5e3da46b/pkg/selfupdate/selfupdate.go
// https://github.com/hyperhq/hypercli/blob/302a6b530148f6a777cd6b8772f706ab5e3da46b/hyper/hyper.go
// https://github.com/inconshreveable/go-update
//
func (e *UpdateAction) Start() error {
	latestReleaseReq, err := http.Get("https://api.github.com/repos/AnalogJ/drawbridge/releases/latest")
	if err != nil {
		return err
	}
	defer latestReleaseReq.Body.Close()

	respBodyJson, err := ioutil.ReadAll(latestReleaseReq.Body)
	if err != nil {
		return err
	}

	//parse json
	releaseInfo := GithubReleaseInfo{}
	json.Unmarshal(respBodyJson, &releaseInfo)
	releaseVersion := fmt.Sprintf("v%v", releaseInfo.TagName)

	//compare the current version to the destination version
	currentTimestamp, err := e.currentBinaryTimestamp()
	if err != nil{
		return err
	}

	fmt.Printf("Current: v%v [%v]. Available: %v [%v]", e.currentBinaryVersion(),currentTimestamp.Format("2006-01-02") , releaseVersion, releaseInfo.PublishedAt.Format("2006-01-02") )

	if releaseVersion == e.currentBinaryVersion(){
		return nil
	}

	//see if theres a binary for this OS/Arch
	assetUrl := ""
	requiredOsArch := fmt.Sprintf("drawbridge-%v-%v", runtime.GOOS, runtime.GOARCH)
	for _, asset := range releaseInfo.Assets{
		if asset.Name == requiredOsArch {
			assetUrl = asset.Url
		}
	}

	if len(assetUrl) == 0 {
		return errors.UpdateBinaryOsArchMissingError(fmt.Sprintf("Cannot find a drawbridge binary for OS/Arch: %v", requiredOsArch))
	}


	//TODO: ask user if we should update.

	releaseBinaryReq, err := http.Get(assetUrl)
	if err != nil {
		return err
	}
	defer releaseBinaryReq.Body.Close()


	err = update.Apply(releaseBinaryReq.Body, update.Options{})
	if err != nil {
		// error handling
		return err
	}


	return nil
}

func (e *UpdateAction) currentBinaryTimestamp() (time.Time, error) {
	execPath, err := os.Executable()
	if err != nil {
		return time.Time{}, err
	}
	info, err := os.Stat(execPath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), err
}

func (e *UpdateAction) currentBinaryVersion() string {
	return version.VERSION
}