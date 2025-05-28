package azurekeyvault

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/azure-keyvault"
	azcommon "github.com/usual2970/certimate/internal/pkg/sdk3rd/azure/common"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

type DeployerConfig struct {
	// Azure TenantId。
	TenantId string `json:"tenantId"`
	// Azure ClientId。
	ClientId string `json:"clientId"`
	// Azure ClientSecret。
	ClientSecret string `json:"clientSecret"`
	// Azure 主权云环境。
	CloudName string `json:"cloudName,omitempty"`
	// Key Vault 名称。
	KeyVaultName string `json:"keyvaultName"`
	// Key Vault 证书名称。
	// 选填。零值时表示新建证书；否则表示更新证书。
	CertificateName string `json:"certificateName,omitempty"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *azcertificates.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.TenantId, config.ClientId, config.ClientSecret, config.CloudName, config.KeyVaultName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdk client: %w", err)
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		TenantId:     config.TenantId,
		ClientId:     config.ClientId,
		ClientSecret: config.ClientSecret,
		CloudName:    config.CloudName,
		KeyVaultName: config.KeyVaultName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ssl uploader: %w", err)
	}

	return &DeployerProvider{
		config:      config,
		logger:      slog.Default(),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
	d.sslUploader.WithLogger(logger)
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*deployer.DeployResult, error) {
	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 转换证书格式
	certPFX, err := certutil.TransformCertificateFromPEMToPFX(certPEM, privkeyPEM, "")
	if err != nil {
		return nil, fmt.Errorf("failed to transform certificate from PEM to PFX: %w", err)
	}

	if d.config.CertificateName == "" {
		// 上传证书到 KeyVault
		upres, err := d.sslUploader.Upload(ctx, certPEM, privkeyPEM)
		if err != nil {
			return nil, fmt.Errorf("failed to upload certificate file: %w", err)
		} else {
			d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
		}
	} else {
		// 获取证书
		// REF: https://learn.microsoft.com/en-us/rest/api/keyvault/certificates/get-certificate/get-certificate
		getCertificateResp, err := d.sdkClient.GetCertificate(context.TODO(), d.config.CertificateName, "", nil)
		d.logger.Debug("sdk request 'keyvault.GetCertificate'", slog.String("request.certificateName", d.config.CertificateName), slog.Any("response", getCertificateResp))
		if err != nil {
			var respErr *azcore.ResponseError
			if !errors.As(err, &respErr) || (respErr.ErrorCode != "ResourceNotFound" && respErr.ErrorCode != "CertificateNotFound") {
				return nil, fmt.Errorf("failed to execute sdk request 'keyvault.GetCertificate': %w", err)
			}
		} else {
			oldCertX509, err := x509.ParseCertificate(getCertificateResp.CER)
			if err == nil {
				if certutil.EqualCertificate(certX509, oldCertX509) {
					return &deployer.DeployResult{}, nil
				}
			}
		}

		// 导入证书
		// REF: https://learn.microsoft.com/en-us/rest/api/keyvault/certificates/import-certificate/import-certificate
		importCertificateParams := azcertificates.ImportCertificateParameters{
			Base64EncodedCertificate: to.Ptr(base64.StdEncoding.EncodeToString(certPFX)),
			CertificatePolicy: &azcertificates.CertificatePolicy{
				SecretProperties: &azcertificates.SecretProperties{
					ContentType: to.Ptr("application/x-pkcs12"),
				},
			},
			Tags: map[string]*string{
				"certimate/cert-cn": to.Ptr(certX509.Subject.CommonName),
				"certimate/cert-sn": to.Ptr(certX509.SerialNumber.Text(16)),
			},
		}
		importCertificateResp, err := d.sdkClient.ImportCertificate(context.TODO(), d.config.CertificateName, importCertificateParams, nil)
		d.logger.Debug("sdk request 'keyvault.ImportCertificate'", slog.String("request.certificateName", d.config.CertificateName), slog.Any("request.parameters", importCertificateParams), slog.Any("response", importCertificateResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'keyvault.ImportCertificate': %w", err)
		}
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(tenantId, clientId, clientSecret, cloudName, keyvaultName string) (*azcertificates.Client, error) {
	env, err := azcommon.GetCloudEnvironmentConfiguration(cloudName)
	if err != nil {
		return nil, err
	}
	clientOptions := azcore.ClientOptions{Cloud: env}

	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret,
		&azidentity.ClientSecretCredentialOptions{ClientOptions: clientOptions})
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("https://%s.vault.azure.net", keyvaultName)
	if azcommon.IsEnvironmentGovernment(cloudName) {
		endpoint = fmt.Sprintf("https://%s.vault.usgovcloudapi.net", keyvaultName)
	} else if azcommon.IsEnvironmentChina(cloudName) {
		endpoint = fmt.Sprintf("https://%s.vault.azure.cn", keyvaultName)
	}

	client, err := azcertificates.NewClient(endpoint, credential, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
