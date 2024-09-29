package applicant

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/pkg/sftp"
	sshPkg "golang.org/x/crypto/ssh"
)

type ssh struct {
	option *ApplyOption
	infos  []string
}

// Apply implements Applicant.
func (s *ssh) Apply() (*Certificate, error) {
	return applyWithFile(s.option, s)
}

type sshAccess struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Key      string `json:"key"`
	Port     string `json:"port"`
}

func NewSSHApplicant(option *ApplyOption) (Applicant, error) {
	return &ssh{
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (s *ssh) upload(client *sshPkg.Client, content, targetPath string) error {

	sftpCli, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create sftp client: %w", err)
	}
	defer sftpCli.Close()

	if err := sftpCli.MkdirAll(path.Dir(targetPath)); err != nil {
		return fmt.Errorf("failed to create remote directory: %w", err)
	}

	file, err := sftpCli.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
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

func toStr(tag string, data any) string {
	if data == nil {
		return tag
	}
	byts, _ := json.Marshal(data)
	return tag + "：" + string(byts)
}

// Present starts a web server and makes the token available at `ChallengePath(token)` for web requests.
func (s *ssh) Present(domain, token, keyAuth string) error {
	access := &sshAccess{}
	if err := json.Unmarshal([]byte(s.option.Access), access); err != nil {
		return err
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

	if value, ok := s.option.Extra["challengeFilePath"]; ok {
		// 上传验证文件
		if err := s.upload(client, keyAuth, fmt.Sprintf("%s/.well-known/acme-challenge/%s", value, token)); err != nil {
			return fmt.Errorf("failed to upload verify file: %w", err)
		}
	} else {
		return fmt.Errorf("verify file path undefined")
	}

	return nil
}

// CleanUp closes the HTTP server and removes the token from `ChallengePath(token)`.
func (s *ssh) CleanUp(domain, token, keyAuth string) error {

	return nil
}
