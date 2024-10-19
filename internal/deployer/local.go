package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"certimate/internal/domain"
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

	preCommand := getDeployString(d.option.DeployConfig, "preCommand")

	if preCommand != "" {
		if err := execCmd(preCommand); err != nil {
			return fmt.Errorf("执行前置命令失败: %w", err)
		}
	}

	// 复制证书文件
	if err := copyFile(getDeployString(d.option.DeployConfig, "certPath"), d.option.Certificate.Certificate); err != nil {
		return fmt.Errorf("复制证书失败: %w", err)
	}

	// 复制私钥文件
	if err := copyFile(getDeployString(d.option.DeployConfig, "keyPath"), d.option.Certificate.PrivateKey); err != nil {
		return fmt.Errorf("复制私钥失败: %w", err)
	}

	// 执行命令
	if err := execCmd(getDeployString(d.option.DeployConfig, "command")); err != nil {
		return fmt.Errorf("执行命令失败: %w", err)
	}

	return nil
}

func execCmd(command string) error {
	// 执行命令
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("执行命令失败: %w", err)
	}

	return nil
}

func copyFile(path string, content string) error {
	dir := filepath.Dir(path)

	// 如果目录不存在，创建目录
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建或打开文件
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 写入内容到文件
	_, err = file.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}
