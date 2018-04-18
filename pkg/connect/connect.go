package connect

import (
	"drawbridge/pkg/config"
	"drawbridge/pkg/utils"
	"syscall"
	"fmt"
	"path/filepath"
)

type ConnectEngine struct {
	Config       config.Interface
}

func (e *ConnectEngine) Start(answerData map[string]interface{}) error {

	//"-c", "command1; command2; command3; ..."

	tmplData, err := e.Config.GetActiveConfigTemplate()
	if err != nil {
		return nil
	}
	tmplConfigName, err := utils.PopulateTemplate(tmplData.FilePath, answerData)
	if err != nil {
		return nil
	}

	//Print the lines we're running.
	//Check that the bastion host is accessible.

	return syscall.Exec("/bin/bash", []string{"-c",
		fmt.Sprintf("ssh-add %v; ssh bastion -F %v;",
			filepath.Join(e.Config.GetString("options.pem_dir"), e.Config.GetString("options.pem_filename")),
			filepath.Join(e.Config.GetString("options.config_dir"), tmplConfigName)),
	}, []string{});
}
