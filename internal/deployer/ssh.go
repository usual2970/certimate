package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/usual2970/certimate/internal/domain"
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
	client, err := d.createSshClient(access)
	if err != nil {
		return err
	}
	defer client.Close()

	d.infos = append(d.infos, toStr("SSH 连接成功", nil))

	// 执行前置命令
	preCommand := d.option.DeployConfig.GetConfigAsString("preCommand")
	if preCommand != "" {
		stdout, stderr, err := d.sshExecCommand(client, preCommand)
		if err != nil {
			return fmt.Errorf("failed to run pre-command: %w, stdout: %s, stderr: %s", err, stdout, stderr)
		}

		d.infos = append(d.infos, toStr("SSH 执行前置命令成功", stdout))
	}

	// 上传证书和私钥文件
	switch d.option.DeployConfig.GetConfigOrDefaultAsString("format", "pem") {
	case "pem":
		if err := d.writeSftpFileString(client, d.option.DeployConfig.GetConfigAsString("certPath"), d.option.Certificate.Certificate); err != nil {
			return fmt.Errorf("failed to upload certificate file: %w", err)
		}

		d.infos = append(d.infos, toStr("SSH 上传证书成功", nil))

		if err := d.writeSftpFileString(client, d.option.DeployConfig.GetConfigAsString("keyPath"), d.option.Certificate.PrivateKey); err != nil {
			return fmt.Errorf("failed to upload private key file: %w", err)
		}

		d.infos = append(d.infos, toStr("SSH 上传私钥成功", nil))

	case "pfx":
		pfxData, err := convertPemToPfx(d.option.Certificate.Certificate, d.option.Certificate.PrivateKey, d.option.DeployConfig.GetConfigAsString("pfxPassword"))
		if err != nil {
			return fmt.Errorf("failed to convert pem to pfx %w", err)
		}

		if err := d.writeSftpFile(client, d.option.DeployConfig.GetConfigAsString("certPath"), pfxData); err != nil {
			return fmt.Errorf("failed to upload certificate file: %w", err)
		}

		d.infos = append(d.infos, toStr("SSH 上传证书成功", nil))
	}

	// 执行命令
	command := d.option.DeployConfig.GetConfigAsString("command")
	if command != "" {
		stdout, stderr, err := d.sshExecCommand(client, command)
		if err != nil {
			return fmt.Errorf("failed to run command: %w, stdout: %s, stderr: %s", err, stdout, stderr)
		}

		d.infos = append(d.infos, toStr("SSH 执行命令成功", stdout))
	}

	return nil
}

func (d *SSHDeployer) createSshClient(access *domain.SSHAccess) (*ssh.Client, error) {
	var authMethod ssh.AuthMethod

	if access.Key != "" {
		var signer ssh.Signer
		var err error

		if access.KeyPassphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(access.Key), []byte(access.KeyPassphrase))
		} else {
			signer, err = ssh.ParsePrivateKey([]byte(access.Key))
		}

		if err != nil {
			return nil, err
		}
		authMethod = ssh.PublicKeys(signer)
	} else {
		authMethod = ssh.Password(access.Password)
	}

	return ssh.Dial("tcp", fmt.Sprintf("%s:%s", access.Host, access.Port), &ssh.ClientConfig{
		User: access.Username,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
}

func (d *SSHDeployer) sshExecCommand(client *ssh.Client, command string) (string, string, error) {
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

func (d *SSHDeployer) writeSftpFileString(client *ssh.Client, path string, content string) error {
	return d.writeSftpFile(client, path, []byte(content))
}

func (d *SSHDeployer) writeSftpFile(client *ssh.Client, path string, data []byte) error {
	sftpCli, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create sftp client: %w", err)
	}
	defer sftpCli.Close()

	if err := sftpCli.MkdirAll(filepath.Dir(path)); err != nil {
		return fmt.Errorf("failed to create remote directory: %w", err)
	}

	file, err := sftpCli.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to remote file: %w", err)
	}

	return nil
}
