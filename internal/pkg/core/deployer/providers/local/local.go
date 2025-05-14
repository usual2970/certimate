package local

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"runtime"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
	fileutil "github.com/usual2970/certimate/internal/pkg/utils/file"
)

type DeployerConfig struct {
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
	// 输出服务器证书文件路径。
	// 选填。
	OutputServerCertPath string `json:"outputServerCertPath,omitempty"`
	// 输出中间证书文件路径。
	// 选填。
	OutputIntermediaCertPath string `json:"outputIntermediaCertPath,omitempty"`
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

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	// 提取服务器证书和中间证书
	serverCertPEM, intermediaCertPEM, err := certutil.ExtractCertificatesFromPEM(certPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to extract certs: %w", err)
	}

	// 执行前置命令
	if d.config.PreCommand != "" {
		stdout, stderr, err := execCommand(d.config.ShellEnv, d.config.PreCommand)
		d.logger.Debug("run pre-command", slog.String("stdout", stdout), slog.String("stderr", stderr))
		if err != nil {
			return nil, fmt.Errorf("failed to execute pre-command (stdout: %s, stderr: %s): %w ", stdout, stderr, err)
		}
	}

	// 写入证书和私钥文件
	switch d.config.OutputFormat {
	case OUTPUT_FORMAT_PEM:
		if err := fileutil.WriteString(d.config.OutputCertPath, certPEM); err != nil {
			return nil, fmt.Errorf("failed to save certificate file: %w", err)
		}
		d.logger.Info("ssl certificate file saved", slog.String("path", d.config.OutputCertPath))

		if d.config.OutputServerCertPath != "" {
			if err := fileutil.WriteString(d.config.OutputServerCertPath, serverCertPEM); err != nil {
				return nil, fmt.Errorf("failed to save server certificate file: %w", err)
			}
			d.logger.Info("ssl server certificate file saved", slog.String("path", d.config.OutputServerCertPath))
		}

		if d.config.OutputIntermediaCertPath != "" {
			if err := fileutil.WriteString(d.config.OutputIntermediaCertPath, intermediaCertPEM); err != nil {
				return nil, fmt.Errorf("failed to save intermedia certificate file: %w", err)
			}
			d.logger.Info("ssl intermedia certificate file saved", slog.String("path", d.config.OutputIntermediaCertPath))
		}

		if err := fileutil.WriteString(d.config.OutputKeyPath, privkeyPEM); err != nil {
			return nil, fmt.Errorf("failed to save private key file: %w", err)
		}
		d.logger.Info("ssl private key file saved", slog.String("path", d.config.OutputKeyPath))

	case OUTPUT_FORMAT_PFX:
		pfxData, err := certutil.TransformCertificateFromPEMToPFX(certPEM, privkeyPEM, d.config.PfxPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to transform certificate to PFX: %w", err)
		}
		d.logger.Info("ssl certificate transformed to pfx")

		if err := fileutil.Write(d.config.OutputCertPath, pfxData); err != nil {
			return nil, fmt.Errorf("failed to save certificate file: %w", err)
		}
		d.logger.Info("ssl certificate file saved", slog.String("path", d.config.OutputCertPath))

	case OUTPUT_FORMAT_JKS:
		jksData, err := certutil.TransformCertificateFromPEMToJKS(certPEM, privkeyPEM, d.config.JksAlias, d.config.JksKeypass, d.config.JksStorepass)
		if err != nil {
			return nil, fmt.Errorf("failed to transform certificate to JKS: %w", err)
		}
		d.logger.Info("ssl certificate transformed to jks")

		if err := fileutil.Write(d.config.OutputCertPath, jksData); err != nil {
			return nil, fmt.Errorf("failed to save certificate file: %w", err)
		}
		d.logger.Info("ssl certificate file saved", slog.String("path", d.config.OutputCertPath))

	default:
		return nil, fmt.Errorf("unsupported output format '%s'", d.config.OutputFormat)
	}

	// 执行后置命令
	if d.config.PostCommand != "" {
		stdout, stderr, err := execCommand(d.config.ShellEnv, d.config.PostCommand)
		d.logger.Debug("run post-command", slog.String("stdout", stdout), slog.String("stderr", stderr))
		if err != nil {
			return nil, fmt.Errorf("failed to execute post-command (stdout: %s, stderr: %s): %w ", stdout, stderr, err)
		}
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

	case ShellEnvType(""):
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}

	default:
		return "", "", fmt.Errorf("unsupported shell env '%s'", shellEnv)
	}

	stdoutBuf := bytes.NewBuffer(nil)
	cmd.Stdout = stdoutBuf
	stderrBuf := bytes.NewBuffer(nil)
	cmd.Stderr = stderrBuf
	err := cmd.Run()
	if err != nil {
		return stdoutBuf.String(), stderrBuf.String(), fmt.Errorf("failed to execute command: %w", err)
	}

	return stdoutBuf.String(), stderrBuf.String(), nil
}
