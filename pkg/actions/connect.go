package actions

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/analogj/drawbridge/pkg/config"
	"github.com/analogj/drawbridge/pkg/errors"
	"github.com/analogj/drawbridge/pkg/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type ConnectAction struct {
	Config config.Interface
}

func (e *ConnectAction) Start(answerData map[string]interface{}, destHostname string) error {
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

	//https://gobyexample.com/execing-processes
	//https://groob.io/posts/golang-execve/

	fmt.Println("Opening ssh tunnel")
	sshBin, lookErr := exec.LookPath("ssh")
	if lookErr != nil {
		return errors.DependencyMissingError("ssh is missing")
	}

	configHost := "bastion"
	if len(destHostname) > 0 {
		configHost = fmt.Sprintf("%v+%v", configHost, destHostname)
	}
	args := []string{"ssh", configHost, "-F", tmplConfigFilepath}

	return syscall.Exec(sshBin, args, os.Environ())
}

func (e *ConnectAction) SshAgentAddPemKey(pemFilepath string) error {
	//first lets ensure that the pemFilepath exists
	if !utils.FileExists(pemFilepath) {
		return errors.PemKeyMissingError(fmt.Sprintf("No pem file exists at %v", pemFilepath))
	}

	//ensure that the ssh-agent is available on this machine.
	_, err := exec.LookPath("ssh-agent")
	if err != nil {
		return errors.DependencyMissingError("ssh-agent is missing")
	}

	//read the pem file data
	keyData, err := ioutil.ReadFile(pemFilepath)
	if err != nil {
		return err
	}

	//TODO: check if this pemfile is already added to the ssh-agent

	//decode the ssh pem key (and handle encypted/passphrase protected keys)
	//https://stackoverflow.com/questions/42105432/how-to-use-an-encrypted-private-key-with-golang-ssh
	block, _ := pem.Decode(keyData)

	//https://github.com/golang/crypto/blob/master/ssh/keys.go

	fmt.Printf("Adding PEM key (%v) to ssh-agent\n", pemFilepath)

	var privateKeyData interface{}
	if x509.IsEncryptedPEMBlock(block) {
		//inform the user that the key is encrypted.

		passphrase, err := utils.StdinQueryPassword(fmt.Sprintf("The key at %v is encrypted and requires a passphrase. Please enter it below:", pemFilepath))
		if err != nil {
			return err
		}

		privateKeyData, err = ssh.ParseRawPrivateKeyWithPassphrase(keyData, []byte(passphrase))
	} else {
		privateKeyData, err = ssh.ParseRawPrivateKey(keyData)
	}

	// register the privatekey with ssh-agent

	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return err
	}
	agentClient := agent.NewClient(conn)

	err = agentClient.Add(agent.AddedKey{
		PrivateKey:   privateKeyData,
		Comment:      fmt.Sprintf("(drawbridge) - %v", pemFilepath),
		LifetimeSecs: 3600, //for safety we should limit this key's use for 1h
	})

	return err
}
