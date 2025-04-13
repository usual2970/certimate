package wangsucdnpro

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/utils/certutil"
	wangsucdn "github.com/usual2970/certimate/internal/pkg/vendors/wangsu-sdk/cdn"
)

type DeployerConfig struct {
	// 网宿云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 网宿云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 网宿云环境。
	Environment string `json:"environment"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
	// 证书 ID。
	// 选填。
	CertificateId string `json:"certificateId,omitempty"`
	// Webhook ID。
	// 选填。
	WebhookId string `json:"webhookId,omitempty"`
}

type DeployerProvider struct {
	config    *DeployerConfig
	logger    *slog.Logger
	sdkClient *wangsucdn.Client
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
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

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 解析证书内容
	certX509, err := certutil.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 查询已部署加速域名的详情
	getHostnameDetailResp, err := d.sdkClient.GetHostnameDetail(d.config.Domain)
	d.logger.Debug("sdk request 'cdn.GetHostnameDetail'", slog.String("hostname", d.config.Domain), slog.Any("response", getHostnameDetailResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.GetHostnameDetail'")
	}

	// 生成网宿云证书参数
	encryptedPrivateKey, err := encryptPrivateKey(privkeyPem, d.config.AccessKeySecret, time.Now().Unix())
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to encrypt private key")
	}
	certificateNewVersionInfo := &wangsucdn.CertificateVersion{
		PrivateKey:  tea.String(encryptedPrivateKey),
		Certificate: tea.String(certPem),
		IdentificationInfo: &wangsucdn.CertificateVersionIdentificationInfo{
			CommonName:              tea.String(certX509.Subject.CommonName),
			SubjectAlternativeNames: &certX509.DNSNames,
		},
	}

	// 网宿云证书 URL 中包含证书 ID 及版本号
	// 格式：
	//    http://open.chinanetcenter.com/cdn/certificates/5dca2205f9e9cc0001df7b33
	//    http://open.chinanetcenter.com/cdn/certificates/329f12c1fe6708c23c31e91f/versions/5
	var wangsuCertUrl string
	var wangsuCertId, wangsuCertVer string

	// 如果原证书 ID 为空，则创建证书；否则更新证书。
	timestamp := time.Now().Unix()
	if d.config.CertificateId == "" {
		// 创建证书
		createCertificateReq := &wangsucdn.CreateCertificateRequest{
			Timestamp:  timestamp,
			Name:       tea.String(fmt.Sprintf("certimate_%d", time.Now().UnixMilli())),
			AutoRenew:  tea.String("Off"),
			NewVersion: certificateNewVersionInfo,
		}
		createCertificateResp, err := d.sdkClient.CreateCertificate(createCertificateReq)
		d.logger.Debug("sdk request 'cdn.CreateCertificate'", slog.Any("request", createCertificateReq), slog.Any("response", createCertificateResp))
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.CreateCertificate'")
		}

		wangsuCertUrl = createCertificateResp.CertificateUrl
		d.logger.Info("ssl certificate uploaded", slog.Any("certUrl", wangsuCertUrl))

		wangsuCertIdMatches := regexp.MustCompile(`/certificates/([a-zA-Z0-9-]+)`).FindStringSubmatch(wangsuCertUrl)
		if len(wangsuCertIdMatches) > 1 {
			wangsuCertId = wangsuCertIdMatches[1]
		}

		wangsuCertVer = "1"
	} else {
		// 更新证书
		updateCertificateReq := &wangsucdn.UpdateCertificateRequest{
			Timestamp:  timestamp,
			Name:       tea.String(fmt.Sprintf("certimate_%d", time.Now().UnixMilli())),
			AutoRenew:  tea.String("Off"),
			NewVersion: certificateNewVersionInfo,
		}
		updateCertificateResp, err := d.sdkClient.UpdateCertificate(d.config.CertificateId, updateCertificateReq)
		d.logger.Debug("sdk request 'cdn.CreateCertificate'", slog.Any("certificateId", d.config.CertificateId), slog.Any("request", updateCertificateReq), slog.Any("response", updateCertificateResp))
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.UpdateCertificate'")
		}

		wangsuCertUrl = updateCertificateResp.CertificateUrl
		d.logger.Info("ssl certificate uploaded", slog.Any("certUrl", wangsuCertUrl))

		wangsuCertIdMatches := regexp.MustCompile(`/certificates/([a-zA-Z0-9-]+)`).FindStringSubmatch(wangsuCertUrl)
		if len(wangsuCertIdMatches) > 1 {
			wangsuCertId = wangsuCertIdMatches[1]
		}

		wangsuCertVerMatches := regexp.MustCompile(`/versions/(\d+)`).FindStringSubmatch(wangsuCertUrl)
		if len(wangsuCertVerMatches) > 1 {
			wangsuCertVer = wangsuCertVerMatches[1]
		}
	}

	// 创建部署任务
	// REF: https://www.wangsu.com/document/api-doc/27034
	createDeploymentTaskReq := &wangsucdn.CreateDeploymentTaskRequest{
		Name:   tea.String(fmt.Sprintf("certimate_%d", time.Now().UnixMilli())),
		Target: tea.String(d.config.Environment),
		Actions: &[]wangsucdn.DeploymentTaskAction{
			{
				Action:        tea.String("deploy_cert"),
				CertificateId: tea.String(wangsuCertId),
				Version:       tea.String(wangsuCertVer),
			},
		},
	}
	if d.config.WebhookId != "" {
		createDeploymentTaskReq.Webhook = tea.String(d.config.WebhookId)
	}
	createDeploymentTaskResp, err := d.sdkClient.CreateDeploymentTask(createDeploymentTaskReq)
	d.logger.Debug("sdk request 'cdn.CreateCertificate'", slog.Any("request", createDeploymentTaskReq), slog.Any("response", createDeploymentTaskResp))
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.CreateDeploymentTask'")
	}

	// 循环获取部署任务详细信息，等待任务状态变更
	// REF: https://www.wangsu.com/document/api-doc/27038
	var wangsuTaskId string
	wangsuTaskMatches := regexp.MustCompile(`/deploymentTasks/([a-zA-Z0-9-]+)`).FindStringSubmatch(wangsuCertUrl)
	if len(wangsuTaskMatches) > 1 {
		wangsuTaskId = wangsuTaskMatches[1]
	}
	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		getDeploymentTaskDetailResp, err := d.sdkClient.GetDeploymentTaskDetail(wangsuTaskId)
		d.logger.Debug("sdk request 'cdn.GetDeploymentTaskDetail'", slog.Any("taskId", wangsuTaskId), slog.Any("response", getDeploymentTaskDetailResp))
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to execute sdk request 'cdn.GetDeploymentTaskDetail'")
		}

		if getDeploymentTaskDetailResp.Status == "failed" {
			return nil, errors.New("unexpected deployment task status")
		} else if getDeploymentTaskDetailResp.Status == "succeeded" {
			break
		}

		d.logger.Info("waiting for deployment task completion ...")
		time.Sleep(time.Second * 15)
	}

	return &deployer.DeployResult{}, nil
}

func createSdkClient(accessKeyId, accessKeySecret string) (*wangsucdn.Client, error) {
	if accessKeyId == "" {
		return nil, errors.New("invalid wangsu access key id")
	}

	if accessKeySecret == "" {
		return nil, errors.New("invalid wangsu access key secret")
	}

	return wangsucdn.NewClient(accessKeyId, accessKeySecret), nil
}

func encryptPrivateKey(privkeyPem string, secretKey string, timestamp int64) (string, error) {
	date := time.Unix(timestamp, 0).UTC()
	dateStr := date.Format("Mon, 02 Jan 2006 15:04:05 GMT")

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(dateStr))
	aesivkey := mac.Sum(nil)
	aesivkeyHex := hex.EncodeToString(aesivkey)

	if len(aesivkeyHex) != 64 {
		return "", fmt.Errorf("invalid hmac length: %d", len(aesivkeyHex))
	}
	ivHex := aesivkeyHex[:32]
	keyHex := aesivkeyHex[32:64]

	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode iv: %w", err)
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainBytes := []byte(privkeyPem)
	padlen := aes.BlockSize - len(plainBytes)%aes.BlockSize
	if padlen > 0 {
		paddata := bytes.Repeat([]byte{byte(padlen)}, padlen)
		plainBytes = append(plainBytes, paddata...)
	}

	encBytes := make([]byte, len(plainBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encBytes, plainBytes)

	return base64.StdEncoding.EncodeToString(encBytes), nil
}
