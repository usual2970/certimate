package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	xpath "path"
	"strings"

	"github.com/pkg/sftp"
	sshPkg "golang.org/x/crypto/ssh"
)

type ssh struct {
	option *DeployerOption
	infos  []string
}

type sshAccess struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Key      string `json:"key"`
	Port     string `json:"port"`
	Command  string `json:"command"`
	CertPath string `json:"certPath"`
	KeyPath  string `json:"keyPath"`
}

func NewSSH(option *DeployerOption) (Deployer, error) {
	return &ssh{
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (a *ssh) GetID() string {
	return fmt.Sprintf("%s-%s", a.option.AceessRecord.GetString("name"), a.option.AceessRecord.Id)
}

func (s *ssh) GetInfo() []string {
	return s.infos
}

func (s *ssh) Deploy(ctx context.Context) error {
	access := &sshAccess{}
	if err := json.Unmarshal([]byte(s.option.Access), access); err != nil {
		return err
	}

	// 将证书路径和命令中的变量替换为实际值
	for k, v := range s.option.Variables {
		key := fmt.Sprintf("${%s}", k)
		access.CertPath = strings.ReplaceAll(access.CertPath, key, v)
		access.KeyPath = strings.ReplaceAll(access.KeyPath, key, v)
		access.Command = strings.ReplaceAll(access.Command, key, v)
	}

	// 连接
	client, err := s.getClient(access)
	if err != nil {
		return err
	}
	defer client.Close()

	s.infos = append(s.infos, toStr("ssh连接成功", nil))

	// 上传
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	s.infos = append(s.infos, toStr("ssh创建session成功", nil))

	// 上传证书
	if err := s.upload(client, s.option.Certificate.Certificate, access.CertPath); err != nil {
		return fmt.Errorf("failed to upload certificate: %w", err)
	}

	s.infos = append(s.infos, toStr("ssh上传证书成功", nil))

	// 上传私钥
	if err := s.upload(client, s.option.Certificate.PrivateKey, access.KeyPath); err != nil {
		return fmt.Errorf("failed to upload private key: %w", err)
	}

	s.infos = append(s.infos, toStr("ssh上传私钥成功", nil))

	// 执行命令
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	var stderrBuf bytes.Buffer
	session.Stderr = &stderrBuf

	if err := session.Run(access.Command); err != nil {
		return fmt.Errorf("failed to run command: %w, stdout: %s, stderr: %s", err, stdoutBuf.String(), stderrBuf.String())
	}

	s.infos = append(s.infos, toStr("ssh执行命令成功", []string{stdoutBuf.String()}))

	return nil
}

func (s *ssh) upload(client *sshPkg.Client, content, path string) error {

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

func (s *ssh) getClient(access *sshAccess) (*sshPkg.Client, error) {

	var authMethod sshPkg.AuthMethod

	if access.Key != "" {
		signer, err := sshPkg.ParsePrivateKey([]byte(access.Key))
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
