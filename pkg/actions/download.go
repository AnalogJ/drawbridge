package actions

import (
	"drawbridge/pkg/config"
	"drawbridge/pkg/errors"
	"drawbridge/pkg/utils"
	"fmt"
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


	tmplData, err := e.Config.GetActiveConfigTemplate()
	if err != nil {
		return nil
	}
	tmplConfigFilepath, err := utils.PopulateTemplate(tmplData.FilePath, answerData)
	if err != nil {
		return nil
	}
	tmplConfigFilepath, err = utils.ExpandPath(filepath.Join(e.Config.GetString("options.config_dir"), tmplConfigFilepath))
	if err != nil {
		return nil
	}

	tmplPemFilepath, err := utils.PopulateTemplate(tmplData.PemFilePath, answerData)
	if err != nil {
		return nil
	}
	tmplPemFilepath, err = utils.ExpandPath(filepath.Join(e.Config.GetString("options.pem_dir"), tmplPemFilepath))
	if err != nil {
		return nil
	}

	//TODO: Print the lines we're running.

	//TODO: Check that the bastion host is accessible.

	err = e.SshAgentAddPemKey(tmplPemFilepath)
	if err != nil {
		return err
	}

	fmt.Println("Begin downloading file through bastion")
	scpBin, lookErr := exec.LookPath("scp")
	if lookErr != nil {
		return errors.DependencyMissingError("scp is missing")
	}

	args := []string{"scp", "-F", tmplConfigFilepath, fmt.Sprintf("bastion+%v:%v", destHostname, remoteFilePath), localFilePath}

	return syscall.Exec(scpBin, args, os.Environ())
}

