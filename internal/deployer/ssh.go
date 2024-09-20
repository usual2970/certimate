package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	xpath "path"
	"strings"
	"sync"

	"github.com/pkg/sftp"
	sshPkg "golang.org/x/crypto/ssh"
)

type ssh struct {
	option *DeployerOption
	infos  []string
}

type sshAccess struct {
	Host     []string `json:"host"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Key      string   `json:"key"`
	Port     string   `json:"port"`
	Command  string   `json:"command"`
	CertPath string   `json:"certPath"`
	KeyPath  string   `json:"keyPath"`
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

	// 替换变量
	for k, v := range s.option.Variables {
		key := fmt.Sprintf("${%s}", k)
		access.CertPath = strings.ReplaceAll(access.CertPath, key, v)
		access.KeyPath = strings.ReplaceAll(access.KeyPath, key, v)
		access.Command = strings.ReplaceAll(access.Command, key, v)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(access.Host))
	for _, host := range access.Host {

		wg.Add(1)

		go func(host string) {
			defer wg.Done()

			// 为每个主机创建 SSH 连接
			client, err := s.getClient(access, host)
			if err != nil {
				errChan <- fmt.Errorf("failed to connect to host %s: %w", host, err)
				return
			}
			defer client.Close()

			s.infos = append(s.infos, toStr(fmt.Sprintf("ssh连接成功: %s", host), nil))

			// 并发上传证书和私钥
			if err := s.upload(client, s.option.Certificate.Certificate, access.CertPath); err != nil {
				errChan <- fmt.Errorf("failed to upload certificate to host %s: %w", host, err)
				return
			}

			s.infos = append(s.infos, toStr(fmt.Sprintf("ssh上传证书成功: %s", host), nil))

			if err := s.upload(client, s.option.Certificate.PrivateKey, access.KeyPath); err != nil {
				errChan <- fmt.Errorf("failed to upload private key to host %s: %w", host, err)
				return
			}

			s.infos = append(s.infos, toStr(fmt.Sprintf("ssh上传私钥成功: %s", host), nil))

			// 执行命令
			session, err := client.NewSession()
			if err != nil {
				errChan <- fmt.Errorf("failed to create session on host %s: %w", host, err)
				return
			}
			defer session.Close()

			var stdoutBuf, stderrBuf bytes.Buffer
			session.Stdout = &stdoutBuf
			session.Stderr = &stderrBuf

			if err := session.Run(access.Command); err != nil {
				errChan <- fmt.Errorf("failed to run command on host %s: %w, stdout: %s, stderr: %s", host, err, stdoutBuf.String(), stderrBuf.String())
				return
			}

			s.infos = append(s.infos, toStr(fmt.Sprintf("ssh执行命令成功: %s", host), []string{stdoutBuf.String()}))

		}(host)
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}

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

func (s *ssh) getClient(access *sshAccess, host string) (*sshPkg.Client, error) {
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

	return sshPkg.Dial("tcp", fmt.Sprintf("%s:%s", host, access.Port), &sshPkg.ClientConfig{
		User: access.Username,
		Auth: []sshPkg.AuthMethod{
			authMethod,
		},
		HostKeyCallback: sshPkg.InsecureIgnoreHostKey(),
	})
}
