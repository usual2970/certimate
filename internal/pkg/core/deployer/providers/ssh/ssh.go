package ssh

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	xerrors "github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/povsister/scp"
	"golang.org/x/crypto/ssh"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/utils/certutil"
)

type DeployerConfig struct {
	// SSH 主机。
	// 零值时默认为 "localhost"。
	SshHost string `json:"sshHost,omitempty"`
	// SSH 端口。
	// 零值时默认为 22。
	SshPort int32 `json:"sshPort,omitempty"`
	// SSH 登录用户名。
	SshUsername string `json:"sshUsername,omitempty"`
	// SSH 登录密码。
	SshPassword string `json:"sshPassword,omitempty"`
	// SSH 登录私钥。
	SshKey string `json:"sshKey,omitempty"`
	// SSH 登录私钥口令。
	SshKeyPassphrase string `json:"sshKeyPassphrase,omitempty"`
	// 是否回退使用 SCP。
	UseSCP bool `json:"useSCP,omitempty"`
	// 前置命令。
	PreCommand string `json:"preCommand,omitempty"`
	// 后置命令。
	PostCommand string `json:"postCommand,omitempty"`
	// 输出证书格式。
	OutputFormat OutputFormatType `json:"outputFormat,omitempty"`
	// 输出证书文件路径。
	OutputCertPath string `json:"outputCertPath,omitempty"`
	// 输出私钥文件路径。
	OutputKeyPath string `json:"outputKeyPath,omitempty"`
	// PFX 导出密码。
	// 证书格式为 PFX 时必填。
	PfxPassword string `json:"pfxPassword,omitempty"`
	// JKS 别名。
	// 证书格式为 JKS 时必填。
	JksAlias string `json:"jksAlias,omitempty"`
	// JKS 密钥密码。
	// 证书格式为 JKS 时必填。
	JksKeypass string `json:"jksKeypass,omitempty"`
	// JKS 存储密码。
	// 证书格式为 JKS 时必填。
	JksStorepass string `json:"jksStorepass,omitempty"`
}

type DeployerProvider struct {
	config *DeployerConfig
	logger *slog.Logger
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	return &DeployerProvider{
		config: config,
		logger: slog.Default(),
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 连接
	client, err := createSshClient(
		d.config.SshHost,
		d.config.SshPort,
		d.config.SshUsername,
		d.config.SshPassword,
		d.config.SshKey,
		d.config.SshKeyPassphrase,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssh client")
	}
	defer client.Close()

	d.logger.Info("ssh connected")

	// 执行前置命令
	if d.config.PreCommand != "" {
		stdout, stderr, err := execSshCommand(client, d.config.PreCommand)
		d.logger.Debug("run pre-command", slog.String("stdout", stdout), slog.String("stderr", stderr))
		if err != nil {
			return nil, xerrors.Wrapf(err, "failed to execute pre-command: stdout: %s, stderr: %s", stdout, stderr)
		}
	}

	// 上传证书和私钥文件
	switch d.config.OutputFormat {
	case OUTPUT_FORMAT_PEM:
		if err := writeFileString(client, d.config.UseSCP, d.config.OutputCertPath, certPem); err != nil {
			return nil, xerrors.Wrap(err, "failed to upload certificate file")
		}
		d.logger.Info("ssl certificate file uploaded", slog.String("path", d.config.OutputCertPath))

		if err := writeFileString(client, d.config.UseSCP, d.config.OutputKeyPath, privkeyPem); err != nil {
			return nil, xerrors.Wrap(err, "failed to upload private key file")
		}
		d.logger.Info("ssl private key file uploaded", slog.String("path", d.config.OutputKeyPath))

	case OUTPUT_FORMAT_PFX:
		pfxData, err := certutil.TransformCertificateFromPEMToPFX(certPem, privkeyPem, d.config.PfxPassword)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to transform certificate to PFX")
		}
		d.logger.Info("ssl certificate transformed to pfx")

		if err := writeFile(client, d.config.UseSCP, d.config.OutputCertPath, pfxData); err != nil {
			return nil, xerrors.Wrap(err, "failed to upload certificate file")
		}
		d.logger.Info("ssl certificate file uploaded", slog.String("path", d.config.OutputCertPath))

	case OUTPUT_FORMAT_JKS:
		jksData, err := certutil.TransformCertificateFromPEMToJKS(certPem, privkeyPem, d.config.JksAlias, d.config.JksKeypass, d.config.JksStorepass)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to transform certificate to JKS")
		}
		d.logger.Info("ssl certificate transformed to jks")

		if err := writeFile(client, d.config.UseSCP, d.config.OutputCertPath, jksData); err != nil {
			return nil, xerrors.Wrap(err, "failed to upload certificate file")
		}
		d.logger.Info("ssl certificate file uploaded", slog.String("path", d.config.OutputCertPath))

	default:
		return nil, fmt.Errorf("unsupported output format: %s", d.config.OutputFormat)
	}

	// 执行后置命令
	if d.config.PostCommand != "" {
		stdout, stderr, err := execSshCommand(client, d.config.PostCommand)
		d.logger.Debug("run post-command", slog.String("stdout", stdout), slog.String("stderr", stderr))
		if err != nil {
			return nil, xerrors.Wrapf(err, "failed to execute post-command, stdout: %s, stderr: %s", stdout, stderr)
		}
	}

	return &deployer.DeployResult{}, nil
}

func createSshClient(host string, port int32, username string, password string, key string, keyPassphrase string) (*ssh.Client, error) {
	if host == "" {
		host = "localhost"
	}

	if port == 0 {
		port = 22
	}

	var authMethod ssh.AuthMethod
	if key != "" {
		var signer ssh.Signer
		var err error

		if keyPassphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(key), []byte(keyPassphrase))
		} else {
			signer, err = ssh.ParsePrivateKey([]byte(key))
		}

		if err != nil {
			return nil, err
		}
		authMethod = ssh.PublicKeys(signer)
	} else {
		authMethod = ssh.Password(password)
	}

	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
}

func execSshCommand(sshCli *ssh.Client, command string) (string, string, error) {
	session, err := sshCli.NewSession()
	if err != nil {
		return "", "", err
	}
	defer session.Close()

	stdoutBuf := bytes.NewBuffer(nil)
	session.Stdout = stdoutBuf
	stderrBuf := bytes.NewBuffer(nil)
	session.Stderr = stderrBuf
	err = session.Run(command)
	if err != nil {
		return stdoutBuf.String(), stderrBuf.String(), xerrors.Wrap(err, "failed to execute ssh command")
	}

	return stdoutBuf.String(), stderrBuf.String(), nil
}

func writeFileString(sshCli *ssh.Client, useSCP bool, path string, content string) error {
	if useSCP {
		return writeFileStringWithSCP(sshCli, path, content)
	}

	return writeFileStringWithSFTP(sshCli, path, content)
}

func writeFile(sshCli *ssh.Client, useSCP bool, path string, data []byte) error {
	if useSCP {
		return writeFileWithSCP(sshCli, path, data)
	}

	return writeFileWithSFTP(sshCli, path, data)
}

func writeFileStringWithSCP(sshCli *ssh.Client, path string, content string) error {
	return writeFileWithSCP(sshCli, path, []byte(content))
}

func writeFileWithSCP(sshCli *ssh.Client, path string, data []byte) error {
	scpCli, err := scp.NewClientFromExistingSSH(sshCli, &scp.ClientOption{})
	if err != nil {
		return xerrors.Wrap(err, "failed to create scp client")
	}

	reader := bytes.NewReader(data)
	err = scpCli.CopyToRemote(reader, path, &scp.FileTransferOption{})
	if err != nil {
		return xerrors.Wrap(err, "failed to write to remote file")
	}

	return nil
}

func writeFileStringWithSFTP(sshCli *ssh.Client, path string, content string) error {
	return writeFileWithSFTP(sshCli, path, []byte(content))
}

func writeFileWithSFTP(sshCli *ssh.Client, path string, data []byte) error {
	sftpCli, err := sftp.NewClient(sshCli)
	if err != nil {
		return xerrors.Wrap(err, "failed to create sftp client")
	}
	defer sftpCli.Close()

	if err := sftpCli.MkdirAll(filepath.Dir(path)); err != nil {
		return xerrors.Wrap(err, "failed to create remote directory")
	}

	file, err := sftpCli.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return xerrors.Wrap(err, "failed to open remote file")
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return xerrors.Wrap(err, "failed to write to remote file")
	}

	return nil
}
