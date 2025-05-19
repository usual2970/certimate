package lecdn

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	leclientsdkv3 "github.com/usual2970/certimate/internal/pkg/sdk3rd/lecdn/v3/client"
	lemastersdkv3 "github.com/usual2970/certimate/internal/pkg/sdk3rd/lecdn/v3/master"
)

type DeployerConfig struct {
	// LeCDN URL。
	ApiUrl string `json:"apiUrl"`
	// LeCDN 版本。
	// 可取值 "v3"。
	ApiVersion string `json:"apiVersion"`
	// LeCDN 用户角色。
	// 可取值 "client"、"master"。
	ApiRole string `json:"apiRole"`
	// LeCDN 用户名。
	Username string `json:"accessKeyId"`
	// LeCDN 用户密码。
	Password string `json:"accessKey"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 证书 ID。
	// 部署资源类型为 [RESOURCE_TYPE_CERTIFICATE] 时必填。
	CertificateId int64 `json:"certificateId,omitempty"`
	// 客户 ID。
	// 部署资源类型为 [RESOURCE_TYPE_CERTIFICATE] 时选填。
	ClientId int64 `json:"clientId,omitempty"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient interface{}
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

const (
	apiVersionV3 = "v3"

	apiRoleClient = "client"
	apiRoleMaster = "master"
)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.ApiUrl, config.ApiVersion, config.ApiRole, config.Username, config.Password, config.AllowInsecureConnections)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	return &DeployerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
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
	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case RESOURCE_TYPE_CERTIFICATE:
		if err := d.deployToCertificate(ctx, certPEM, privkeyPEM); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type '%s'", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *DeployerProvider) deployToCertificate(ctx context.Context, certPEM string, privkeyPEM string) error {
	if d.config.CertificateId == 0 {
		return errors.New("config `certificateId` is required")
	}

	// 修改证书
	// REF: https://wdk0pwf8ul.feishu.cn/wiki/YE1XwCRIHiLYeKkPupgcXrlgnDd
	switch sdkClient := d.sdkClient.(type) {
	case *leclientsdkv3.Client:
		updateSSLCertReq := &leclientsdkv3.UpdateCertificateRequest{
			Name:        fmt.Sprintf("certimate-%d", time.Now().UnixMilli()),
			Description: "upload from certimate",
			Type:        "upload",
			SSLPEM:      certPEM,
			SSLKey:      privkeyPEM,
			AutoRenewal: false,
		}
		updateSSLCertResp, err := sdkClient.UpdateCertificate(d.config.CertificateId, updateSSLCertReq)
		d.logger.Debug("sdk request 'lecdn.UpdateCertificate'", slog.Int64("certId", d.config.CertificateId), slog.Any("request", updateSSLCertReq), slog.Any("response", updateSSLCertResp))
		if err != nil {
			return fmt.Errorf("failed to execute sdk request 'lecdn.UpdateCertificate': %w", err)
		}

	case *lemastersdkv3.Client:
		updateSSLCertReq := &lemastersdkv3.UpdateCertificateRequest{
			ClientId:    d.config.ClientId,
			Name:        fmt.Sprintf("certimate-%d", time.Now().UnixMilli()),
			Description: "upload from certimate",
			Type:        "upload",
			SSLPEM:      certPEM,
			SSLKey:      privkeyPEM,
			AutoRenewal: false,
		}
		updateSSLCertResp, err := sdkClient.UpdateCertificate(d.config.CertificateId, updateSSLCertReq)
		d.logger.Debug("sdk request 'lecdn.UpdateCertificate'", slog.Int64("certId", d.config.CertificateId), slog.Any("request", updateSSLCertReq), slog.Any("response", updateSSLCertResp))
		if err != nil {
			return fmt.Errorf("failed to execute sdk request 'lecdn.UpdateCertificate': %w", err)
		}

	default:
		panic("sdk client is not implemented")
	}

	return nil
}

func createSdkClient(apiUrl, apiVersion, apiRole, username, password string, skipTlsVerify bool) (interface{}, error) {
	if _, err := url.Parse(apiUrl); err != nil {
		return nil, errors.New("invalid lecdn api url")
	}

	if username == "" {
		return nil, errors.New("invalid lecdn username")
	}

	if password == "" {
		return nil, errors.New("invalid lecdn password")
	}

	if apiVersion == apiVersionV3 && apiRole == apiRoleClient {
		// v3 版客户端
		client := leclientsdkv3.NewClient(apiUrl, username, password)
		if skipTlsVerify {
			client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})
		}

		return client, nil
	} else if apiVersion == apiVersionV3 && apiRole == apiRoleMaster {
		// v3 版主控端
		client := lemastersdkv3.NewClient(apiUrl, username, password)
		if skipTlsVerify {
			client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})
		}

		return client, nil
	}

	return nil, fmt.Errorf("invalid lecdn api version or user role")
}
