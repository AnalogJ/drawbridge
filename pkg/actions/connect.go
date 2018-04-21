package actions

import (
	"crypto/x509"
	"drawbridge/pkg/config"
	"drawbridge/pkg/errors"
	"drawbridge/pkg/utils"
	"encoding/pem"
	"fmt"
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

	err = e.SshAgentAddPemKey(tmplPemFilepath)
	if err != nil {
		return err
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

	fmt.Printf("Adding `%v` PEM key to ssh-agent\n", block.Type)

	var privateKeyData interface{}
	if x509.IsEncryptedPEMBlock(block) {
		//inform the user that the key is encrypted.
		passphrase := utils.StdinQuery(fmt.Sprintf("The key at %v is encrypted and requires a passphrase. Please enter it below:", pemFilepath))
		privateKeyData, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(passphrase))
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
		PrivateKey: privateKeyData,
		Comment:    fmt.Sprintf("(drawbridge) - %v", pemFilepath),
		LifetimeSecs: 3600, //for safety we should limit this key's use for 1h
	})

	return err
}
