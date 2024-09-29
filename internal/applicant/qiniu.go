package applicant

import (
	"certimate/internal/domain"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

const qiniuGateway = "http://api.qiniu.com"

type qiuniu struct {
	option      *ApplyOption
	info        []string
	credentials *auth.Credentials
}

// Apply implements Applicant.
func (s *qiuniu) Apply() (*Certificate, error) {
	return applyWithFile(s.option, s)
}

func NewQiNiuApplicant(option *ApplyOption) (Applicant, error) {
	access := &domain.QiniuAccess{}
	json.Unmarshal([]byte(option.Access), access)

	return &qiuniu{
		option: option,
		info:   make([]string, 0),

		credentials: auth.New(access.AccessKey, access.SecretKey),
	}, nil
}

// Present starts a web server and makes the token available at `ChallengePath(token)` for web requests.
func (s *qiuniu) Present(domain, token, keyAuth string) error {

	if value, ok := s.option.Extra["challengeFilePath"]; ok {
		putPolicy := storage.PutPolicy{
			Scope: value,
		}
		upToken := putPolicy.UploadToken(s.credentials)
		cfg := storage.Config{}
		// 是否使用https域名
		cfg.UseHTTPS = true
		// 上传是否使用CDN上传加速
		cfg.UseCdnDomains = false
		uploader := storage.NewBase64Uploader(&cfg)
		ret := storage.PutRet{}
		// 可选配置
		putExtra := storage.Base64PutExtra{
			Params: map[string]string{
				"x:name": token,
			},
		}
		err := uploader.Put(context.Background(), &ret, upToken,
			fmt.Sprintf(".well-known/acme-challenge/%s", token),
			[]byte(base64.StdEncoding.EncodeToString([]byte(keyAuth))),
			&putExtra)

		// 上传验证文件
		if err != nil {
			return fmt.Errorf("failed to upload verify file: %w", err)
		}
	} else {
		return fmt.Errorf("verify file path undefined")
	}

	return nil
}

// CleanUp closes the HTTP server and removes the token from `ChallengePath(token)`.
func (s *qiuniu) CleanUp(domain, token, keyAuth string) error {

	return nil
}
