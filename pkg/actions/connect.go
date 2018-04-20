package actions

import (
	"drawbridge/pkg/config"
	"drawbridge/pkg/utils"
	"path/filepath"
	"syscall"
	"os"
	"os/exec"
)

type ConnectAction struct {
	Config config.Interface
}

func (e *ConnectAction) Start(answerData map[string]interface{}) error {

	//"-c", "command1; command2; command3; ..."

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

	//TODO: add the ssh/pem key to the ssh-agent (if its running).


	//https://gobyexample.com/execing-processes
	//https://groob.io/posts/golang-execve/


	binary, lookErr := exec.LookPath("ssh")
	if lookErr != nil {
		panic(lookErr)
	}

	args := []string{"ssh", "bastion", "-F", tmplConfigFilepath}

	return syscall.Exec(binary, args, os.Environ())
}
