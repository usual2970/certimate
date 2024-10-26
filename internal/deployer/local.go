package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/utils/fs"
)

type LocalDeployer struct {
	option *DeployerOption
	infos  []string
}

func NewLocalDeployer(option *DeployerOption) (Deployer, error) {
	return &LocalDeployer{
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (d *LocalDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *LocalDeployer) GetInfo() []string {
	return []string{}
}

func (d *LocalDeployer) Deploy(ctx context.Context) error {
	access := &domain.LocalAccess{}
	if err := json.Unmarshal([]byte(d.option.Access), access); err != nil {
		return err
	}

	// 执行前置命令
	preCommand := d.option.DeployConfig.GetConfigAsString("preCommand")
	if preCommand != "" {
		stdout, stderr, err := d.execCommand(preCommand)
		if err != nil {
			return fmt.Errorf("failed to run pre-command: %w, stdout: %s, stderr: %s", err, stdout, stderr)
		}

		d.infos = append(d.infos, toStr("执行前置命令成功", stdout))
	}

	// 写入证书和私钥文件
	switch d.option.DeployConfig.GetConfigOrDefaultAsString("format", "pem") {
	case "pem":
		if err := fs.WriteFileString(d.option.DeployConfig.GetConfigAsString("certPath"), d.option.Certificate.Certificate); err != nil {
			return fmt.Errorf("failed to save certificate file: %w", err)
		}

		d.infos = append(d.infos, toStr("保存证书成功", nil))

		if err := fs.WriteFileString(d.option.DeployConfig.GetConfigAsString("keyPath"), d.option.Certificate.PrivateKey); err != nil {
			return fmt.Errorf("failed to save private key file: %w", err)
		}

		d.infos = append(d.infos, toStr("保存私钥成功", nil))

	case "pfx":
		pfxData, err := convertPemToPfx(d.option.Certificate.Certificate, d.option.Certificate.PrivateKey, d.option.DeployConfig.GetConfigAsString("pfxPassword"))
		if err != nil {
			return fmt.Errorf("failed to convert pem to pfx %w", err)
		}

		if err := fs.WriteFile(d.option.DeployConfig.GetConfigAsString("certPath"), pfxData); err != nil {
			return fmt.Errorf("failed to save certificate file: %w", err)
		}

		d.infos = append(d.infos, toStr("保存证书成功", nil))
	}

	// 执行命令
	command := d.option.DeployConfig.GetConfigAsString("command")
	if command != "" {
		stdout, stderr, err := d.execCommand(command)
		if err != nil {
			return fmt.Errorf("failed to run command: %w, stdout: %s, stderr: %s", err, stdout, stderr)
		}

		d.infos = append(d.infos, toStr("执行命令成功", stdout))
	}

	return nil
}

func (d *LocalDeployer) execCommand(command string) (string, string, error) {
	var cmd *exec.Cmd

	switch d.option.DeployConfig.GetConfigAsString("shell") {
	case "sh":
		cmd = exec.Command("sh", "-c", command)

	case "cmd":
		cmd = exec.Command("cmd", "/C", command)

	case "powershell":
		cmd = exec.Command("powershell", "-Command", command)

	case "":
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}

	default:
		return "", "", fmt.Errorf("unsupported shell")
	}

	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("failed to execute script: %w", err)
	}

	return stdoutBuf.String(), stderrBuf.String(), err
}
