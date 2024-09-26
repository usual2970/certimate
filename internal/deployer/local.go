package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type localAccess struct {
	Command  string `json:"command"`
	CertPath string `json:"certPath"`
	KeyPath  string `json:"keyPath"`
}

type local struct {
	option *DeployerOption
	infos  []string
}

func NewLocal(option *DeployerOption) *local {
	return &local{
		option: option,
		infos:  make([]string, 0),
	}
}

func (l *local) GetID() string {
	return fmt.Sprintf("%s-%s", l.option.AceessRecord.GetString("name"), l.option.AceessRecord.Id)
}

func (l *local) GetInfo() []string {
	return []string{}
}

func (l *local) Deploy(ctx context.Context) error {
	access := &localAccess{}
	if err := json.Unmarshal([]byte(l.option.Access), access); err != nil {
		return err
	}
	// 复制文件
	if err := copyFile(l.option.Certificate.Certificate, access.CertPath); err != nil {
		return fmt.Errorf("复制证书失败: %w", err)
	}

	if err := copyFile(l.option.Certificate.PrivateKey, access.KeyPath); err != nil {
		return fmt.Errorf("复制私钥失败: %w", err)
	}

	// 执行命令

	if err := execCmd(access.Command); err != nil {
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

func copyFile(content string, path string) error {
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
