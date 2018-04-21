package actions

import (
	"drawbridge/pkg/config"
	"drawbridge/pkg/errors"
	"drawbridge/pkg/utils"
	"drawbridge/pkg/version"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/inconshreveable/go-update"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

type UpdateAction struct {
	Config config.Interface
}

type GithubReleaseInfo struct {
	TagName         string    `json:"tag_name"`
	PublishedAt     time.Time `json:"published_at"`
	ReleaseNotesUrl string    `json:"html_url"`
	Assets          []struct {
		Url  string `json:"browser_download_url"`
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
	err = json.Unmarshal(respBodyJson, &releaseInfo)
	if err != nil {
		return err
	}

	//compare the current version to the destination version
	currentTimestamp, err := e.currentBinaryTimestamp()
	if err != nil {
		return err
	}

	fmt.Printf("Current: %v [%v]. Available: %v [%v]\nRelease notes are available here: %v\n",
		e.currentBinaryVersion(),
		currentTimestamp.Format("2006-01-02"),
		releaseInfo.TagName,
		releaseInfo.PublishedAt.Format("2006-01-02"),
		releaseInfo.ReleaseNotesUrl,
	)

	if releaseInfo.TagName == e.currentBinaryVersion() {
		//TODO: return errors.UpdateNotAvailableError("No new version found.")
	}

	//see if theres a binary for this OS/Arch
	assetUrl := ""
	requiredOsArch := fmt.Sprintf("drawbridge-%v-%v", runtime.GOOS, runtime.GOARCH)
	for _, asset := range releaseInfo.Assets {
		if asset.Name == requiredOsArch {
			assetUrl = asset.Url
		}
	}

	if len(assetUrl) == 0 {
		return errors.UpdateBinaryOsArchMissingError(fmt.Sprintf("Cannot find a drawbridge binary for OS/Arch: %v", requiredOsArch))
	}

	//TODO: ask user if we should update.
	stdinResp := utils.StdinQuery(fmt.Sprintf("Are you sure you would like to update drawbridge to %v?\nPlease confirm [true/false]:", releaseInfo.TagName))
	val, err := strconv.ParseBool(stdinResp)
	if err != nil {
		return err
	}
	if !val {
		color.Red("Cancelled delete operation.")
		return nil
	}

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
	return fmt.Sprintf("v%v", version.VERSION)
}
