package actions

import (
	"fmt"
	"github.com/analogj/drawbridge/pkg/config"
	"github.com/analogj/drawbridge/pkg/errors"
	"github.com/analogj/drawbridge/pkg/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type DownloadAction struct {
	ConnectAction
	Config config.Interface
}

func (e *DownloadAction) Start(answerData map[string]interface{}, destHostname string, remoteFilePath string, localFilePath string) error {
	log.Debugf("Answer Data: %v", answerData)

	tmplData, err := e.Config.GetActiveConfigTemplate()
	if err != nil {
		return nil
	}

	tmplConfigFilepath, err := utils.PopulatePathTemplate(filepath.Join(e.Config.GetString("options.config_dir"), tmplData.FilePath), answerData)
	if err != nil {
		return nil
	}

	if tmplData.PemFilePath != "" {
		tmplPemFilepath, err := utils.PopulatePathTemplate(filepath.Join(e.Config.GetString("options.pem_dir"), tmplData.PemFilePath), answerData)
		if err != nil {
			return nil
		}

		//TODO: Print the lines we're running.

		//TODO: Check that the bastion host is accessible.

		err = e.SshAgentAddPemKey(tmplPemFilepath)
		if err != nil {
			return err
		}
	}

	fmt.Println("Begin downloading file through bastion")
	scpBin, lookErr := exec.LookPath("scp")
	if lookErr != nil {
		return errors.DependencyMissingError("scp is missing")
	}

	args := []string{"scp", "-F", tmplConfigFilepath, fmt.Sprintf("%v.in:%v", destHostname, remoteFilePath), localFilePath}

	return syscall.Exec(scpBin, args, os.Environ())
}
