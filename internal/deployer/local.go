package deployer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"

	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/utils/fs"
)

type LocalDeployer struct {
	option *DeployerOption
	infos  []string
}

const (
	certFormatPEM = "pem"
	certFormatPFX = "pfx"
	certFormatJKS = "jks"
)

const (
	shellEnvSh         = "sh"
	shellEnvCmd        = "cmd"
	shellEnvPowershell = "powershell"
)

func NewLocalDeployer(option *DeployerOption) (Deployer, error) {
	return &LocalDeployer{
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (d *LocalDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *LocalDeployer) GetInfos() []string {
	return []string{}
}

func (d *LocalDeployer) Deploy(ctx context.Context) error {
	// 执行前置命令
	preCommand := d.option.DeployConfig.GetConfigAsString("preCommand")
	if preCommand != "" {
		stdout, stderr, err := d.execCommand(preCommand)
		if err != nil {
			return xerrors.Wrapf(err, "failed to run pre-command, stdout: %s, stderr: %s", stdout, stderr)
		}

		d.infos = append(d.infos, toStr("执行前置命令成功", stdout))
	}

	// 写入证书和私钥文件
	switch d.option.DeployConfig.GetConfigOrDefaultAsString("format", certFormatPEM) {
	case certFormatPEM:
		if err := fs.WriteFileString(d.option.DeployConfig.GetConfigAsString("certPath"), d.option.Certificate.Certificate); err != nil {
			return err
		}

		d.infos = append(d.infos, toStr("保存证书成功", nil))

		if err := fs.WriteFileString(d.option.DeployConfig.GetConfigAsString("keyPath"), d.option.Certificate.PrivateKey); err != nil {
			return err
		}

		d.infos = append(d.infos, toStr("保存私钥成功", nil))

	case certFormatPFX:
		pfxData, err := convertPEMToPFX(
			d.option.Certificate.Certificate,
			d.option.Certificate.PrivateKey,
			d.option.DeployConfig.GetConfigAsString("pfxPassword"),
		)
		if err != nil {
			return err
		}

		if err := fs.WriteFile(d.option.DeployConfig.GetConfigAsString("certPath"), pfxData); err != nil {
			return err
		}

		d.infos = append(d.infos, toStr("保存证书成功", nil))

	case certFormatJKS:
		jksData, err := convertPEMToJKS(
			d.option.Certificate.Certificate,
			d.option.Certificate.PrivateKey,
			d.option.DeployConfig.GetConfigAsString("jksAlias"),
			d.option.DeployConfig.GetConfigAsString("jksKeypass"),
			d.option.DeployConfig.GetConfigAsString("jksStorepass"),
		)
		if err != nil {
			return err
		}

		if err := fs.WriteFile(d.option.DeployConfig.GetConfigAsString("certPath"), jksData); err != nil {
			return err
		}

		d.infos = append(d.infos, toStr("保存证书成功", nil))
	}

	// 执行命令
	command := d.option.DeployConfig.GetConfigAsString("command")
	if command != "" {
		stdout, stderr, err := d.execCommand(command)
		if err != nil {
			return xerrors.Wrapf(err, "failed to run command, stdout: %s, stderr: %s", stdout, stderr)
		}

		d.infos = append(d.infos, toStr("执行命令成功", stdout))
	}

	return nil
}

func (d *LocalDeployer) execCommand(command string) (string, string, error) {
	var cmd *exec.Cmd

	switch d.option.DeployConfig.GetConfigAsString("shell") {
	case shellEnvSh:
		cmd = exec.Command("sh", "-c", command)

	case shellEnvCmd:
		cmd = exec.Command("cmd", "/C", command)

	case shellEnvPowershell:
		cmd = exec.Command("powershell", "-Command", command)

	case "":
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}

	default:
		return "", "", errors.New("unsupported shell")
	}

	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		return "", "", xerrors.Wrap(err, "failed to execute shell script")
	}

	return stdoutBuf.String(), stderrBuf.String(), nil
}
