package local

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/utils/certs"
	"github.com/usual2970/certimate/internal/pkg/utils/files"
)

type LocalDeployerConfig struct {
	// Shell 执行环境。
	// 零值时默认根据操作系统决定。
	ShellEnv ShellEnvType `json:"shellEnv,omitempty"`
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

type LocalDeployer struct {
	config *LocalDeployerConfig
	logger logger.Logger
}

var _ deployer.Deployer = (*LocalDeployer)(nil)

func New(config *LocalDeployerConfig) (*LocalDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *LocalDeployerConfig, logger logger.Logger) (*LocalDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	return &LocalDeployer{
		logger: logger,
		config: config,
	}, nil
}

func (d *LocalDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 执行前置命令
	if d.config.PreCommand != "" {
		stdout, stderr, err := execCommand(d.config.ShellEnv, d.config.PreCommand)
		if err != nil {
			return nil, xerrors.Wrapf(err, "failed to run pre-command, stdout: %s, stderr: %s", stdout, stderr)
		}

		d.logger.Logt("pre-command executed", stdout)
	}

	// 写入证书和私钥文件
	switch d.config.OutputFormat {
	case OUTPUT_FORMAT_PEM:
		if err := files.WriteString(d.config.OutputCertPath, certPem); err != nil {
			return nil, xerrors.Wrap(err, "failed to save certificate file")
		}

		d.logger.Logt("certificate file saved")

		if err := files.WriteString(d.config.OutputKeyPath, privkeyPem); err != nil {
			return nil, xerrors.Wrap(err, "failed to save private key file")
		}

		d.logger.Logt("private key file saved")

	case OUTPUT_FORMAT_PFX:
		pfxData, err := certs.TransformCertificateFromPEMToPFX(certPem, privkeyPem, d.config.PfxPassword)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to transform certificate to PFX")
		}

		d.logger.Logt("certificate transformed to PFX")

		if err := files.Write(d.config.OutputCertPath, pfxData); err != nil {
			return nil, xerrors.Wrap(err, "failed to save certificate file")
		}

		d.logger.Logt("certificate file saved")

	case OUTPUT_FORMAT_JKS:
		jksData, err := certs.TransformCertificateFromPEMToJKS(certPem, privkeyPem, d.config.JksAlias, d.config.JksKeypass, d.config.JksStorepass)
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to transform certificate to JKS")
		}

		d.logger.Logt("certificate transformed to JKS")

		if err := files.Write(d.config.OutputCertPath, jksData); err != nil {
			return nil, xerrors.Wrap(err, "failed to save certificate file")
		}

		d.logger.Logt("certificate file uploaded")

	default:
		return nil, fmt.Errorf("unsupported output format: %s", d.config.OutputFormat)
	}

	// 执行后置命令
	if d.config.PostCommand != "" {
		stdout, stderr, err := execCommand(d.config.ShellEnv, d.config.PostCommand)
		if err != nil {
			return nil, xerrors.Wrapf(err, "failed to run command, stdout: %s, stderr: %s", stdout, stderr)
		}

		d.logger.Logt("post-command executed", stdout)
	}

	return &deployer.DeployResult{}, nil
}

func execCommand(shellEnv ShellEnvType, command string) (string, string, error) {
	var cmd *exec.Cmd

	switch shellEnv {
	case SHELL_ENV_SH:
		cmd = exec.Command("sh", "-c", command)

	case SHELL_ENV_CMD:
		cmd = exec.Command("cmd", "/C", command)

	case SHELL_ENV_POWERSHELL:
		cmd = exec.Command("powershell", "-Command", command)

	case "":
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}

	default:
		return "", "", fmt.Errorf("unsupported shell env: %s", shellEnv)
	}

	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		return "", "", xerrors.Wrap(err, "failed to execute shell command")
	}

	return stdoutBuf.String(), stderrBuf.String(), nil
}
