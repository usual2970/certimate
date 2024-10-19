package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	xpath "path"

	"github.com/pkg/sftp"
	sshPkg "golang.org/x/crypto/ssh"

	"certimate/internal/domain"
)

type SSHDeployer struct {
	option *DeployerOption
	infos  []string
}

func NewSSHDeployer(option *DeployerOption) (Deployer, error) {
	return &SSHDeployer{
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (d *SSHDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *SSHDeployer) GetInfo() []string {
	return d.infos
}

func (d *SSHDeployer) Deploy(ctx context.Context) error {
	access := &domain.SSHAccess{}
	if err := json.Unmarshal([]byte(d.option.Access), access); err != nil {
		return err
	}

	// 连接
	client, err := d.createClient(access)
	if err != nil {
		return err
	}
	defer client.Close()

	d.infos = append(d.infos, toStr("ssh连接成功", nil))

	// 执行前置命令
	preCommand := getDeployString(d.option.DeployConfig, "preCommand")
	if preCommand != "" {
		stdout, stderr, err := d.sshExecCommand(client, preCommand)
		if err != nil {
			return fmt.Errorf("failed to run pre-command: %w, stdout: %s, stderr: %s", err, stdout, stderr)
		}
	}

	// 上传证书
	if err := d.upload(client, d.option.Certificate.Certificate, getDeployString(d.option.DeployConfig, "certPath")); err != nil {
		return fmt.Errorf("failed to upload certificate: %w", err)
	}

	d.infos = append(d.infos, toStr("ssh上传证书成功", nil))

	// 上传私钥
	if err := d.upload(client, d.option.Certificate.PrivateKey, getDeployString(d.option.DeployConfig, "keyPath")); err != nil {
		return fmt.Errorf("failed to upload private key: %w", err)
	}

	d.infos = append(d.infos, toStr("ssh上传私钥成功", nil))

	// 执行命令
	stdout, stderr, err := d.sshExecCommand(client, getDeployString(d.option.DeployConfig, "command"))
	if err != nil {
		return fmt.Errorf("failed to run command: %w, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	d.infos = append(d.infos, toStr("ssh执行命令成功", stdout))

	return nil
}

func (d *SSHDeployer) sshExecCommand(client *sshPkg.Client, command string) (string, string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("failed to create ssh session: %w", err)
	}

	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	var stderrBuf bytes.Buffer
	session.Stderr = &stderrBuf
	err = session.Run(command)
	return stdoutBuf.String(), stderrBuf.String(), err
}

func (d *SSHDeployer) upload(client *sshPkg.Client, content, path string) error {
	sftpCli, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create sftp client: %w", err)
	}
	defer sftpCli.Close()

	if err := sftpCli.MkdirAll(xpath.Dir(path)); err != nil {
		return fmt.Errorf("failed to create remote directory: %w", err)
	}

	file, err := sftpCli.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %w", err)
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("failed to write to remote file: %w", err)
	}

	return nil
}

func (d *SSHDeployer) createClient(access *domain.SSHAccess) (*sshPkg.Client, error) {
	var authMethod sshPkg.AuthMethod

	if access.Key != "" {
		var signer sshPkg.Signer
		var err error

		if access.KeyPassphrase != "" {
			signer, err = sshPkg.ParsePrivateKeyWithPassphrase([]byte(access.Key), []byte(access.KeyPassphrase))
		} else {
			signer, err = sshPkg.ParsePrivateKey([]byte(access.Key))
		}

		if err != nil {
			return nil, err
		}
		authMethod = sshPkg.PublicKeys(signer)
	} else {
		authMethod = sshPkg.Password(access.Password)
	}

	return sshPkg.Dial("tcp", fmt.Sprintf("%s:%s", access.Host, access.Port), &sshPkg.ClientConfig{
		User: access.Username,
		Auth: []sshPkg.AuthMethod{
			authMethod,
		},
		HostKeyCallback: sshPkg.InsecureIgnoreHostKey(),
	})
}
