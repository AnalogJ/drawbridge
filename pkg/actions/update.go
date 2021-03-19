package actions

import (
	"encoding/json"
	"fmt"
	"github.com/analogj/drawbridge/pkg/config"
	"github.com/analogj/drawbridge/pkg/errors"
	"github.com/analogj/drawbridge/pkg/utils"
	"github.com/analogj/drawbridge/pkg/version"
	"github.com/fatih/color"
	"github.com/inconshreveable/go-update"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
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

	releaseInfo, err := e.GetLatestReleaseInfo()
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
		return errors.UpdateNotAvailableError("No new version found.")
	}

	//see if theres a binary for this OS/Arch
	assetUrl := ""
	var requiredOsArch string
	if runtime.GOOS == "windows" {
		requiredOsArch = fmt.Sprintf("drawbridge-%v-%v.exe", runtime.GOOS, runtime.GOARCH)
	} else {
		requiredOsArch = fmt.Sprintf("drawbridge-%v-%v", runtime.GOOS, runtime.GOARCH)
	}

	for _, asset := range releaseInfo.Assets {
		if asset.Name == requiredOsArch {
			assetUrl = asset.Url
		}
	}

	if len(assetUrl) == 0 {
		return errors.UpdateBinaryOsArchMissingError(fmt.Sprintf("Cannot find a drawbridge binary for OS/Arch: %v", requiredOsArch))
	}

	val := utils.StdinQueryBoolean(fmt.Sprintf("Are you sure you would like to update drawbridge to %v?\nPlease confirm [yes/no]:", releaseInfo.TagName))

	if !val {
		color.Red("Cancelled update operation.")
		return nil
	}

	color.Yellow("Updating Drawbridge binary. Please wait...")
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

	color.Green("Successfully updated Drawbridge")

	return nil
}

func (e *UpdateAction) GetLatestReleaseInfo() (GithubReleaseInfo, error) {
	latestReleaseReq, err := http.Get("https://api.github.com/repos/AnalogJ/drawbridge/releases/latest")
	if err != nil {
		return GithubReleaseInfo{}, err
	}
	defer latestReleaseReq.Body.Close()

	respBodyJson, err := ioutil.ReadAll(latestReleaseReq.Body)
	if err != nil {
		return GithubReleaseInfo{}, err
	}

	//parse json
	releaseInfo := GithubReleaseInfo{}
	err = json.Unmarshal(respBodyJson, &releaseInfo)
	if err != nil {
		return GithubReleaseInfo{}, err
	}
	return releaseInfo, nil
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
