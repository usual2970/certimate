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
	"strconv"
	"time"

	"github.com/certimate-go/certimate/pkg/core"
	wangsucdn "github.com/certimate-go/certimate/pkg/sdk3rd/wangsu/cdnpro"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
	xtypes "github.com/certimate-go/certimate/pkg/utils/types"
)

type SSLDeployerProviderConfig struct {
	// 网宿云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 网宿云 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 网宿云 API Key。
	ApiKey string `json:"apiKey"`
	// 网宿云环境。
	Environment string `json:"environment"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain"`
	// 证书 ID。
	// 选填。零值时表示新建证书；否则表示更新证书。
	CertificateId string `json:"certificateId,omitempty"`
	// Webhook ID。
	// 选填。
	WebhookId string `json:"webhookId,omitempty"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *wangsucdn.Client
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	return &SSLDeployerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 解析证书内容
	certX509, err := xcert.ParseCertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}

	// 查询已部署加速域名的详情
	getHostnameDetailResp, err := d.sdkClient.GetHostnameDetail(d.config.Domain)
	d.logger.Debug("sdk request 'cdnpro.GetHostnameDetail'", slog.String("hostname", d.config.Domain), slog.Any("response", getHostnameDetailResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdnpro.GetHostnameDetail': %w", err)
	}

	// 生成网宿云证书参数
	encryptedPrivateKey, err := encryptPrivateKey(privkeyPEM, d.config.ApiKey, time.Now().Unix())
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt private key: %w", err)
	}
	certificateNewVersionInfo := &wangsucdn.CertificateVersionInfo{
		PrivateKey:  xtypes.ToPtr(encryptedPrivateKey),
		Certificate: xtypes.ToPtr(certPEM),
		IdentificationInfo: &wangsucdn.CertificateVersionIdentificationInfo{
			CommonName:              xtypes.ToPtr(certX509.Subject.CommonName),
			SubjectAlternativeNames: &certX509.DNSNames,
		},
	}

	// 网宿云证书 URL 中包含证书 ID 及版本号
	// 格式：
	//    http://open.chinanetcenter.com/cdn/certificates/5dca2205f9e9cc0001df7b33
	//    http://open.chinanetcenter.com/cdn/certificates/329f12c1fe6708c23c31e91f/versions/5
	var wangsuCertUrl string
	var wangsuCertId string
	var wangsuCertVer int32

	// 如果原证书 ID 为空，则创建证书；否则更新证书。
	timestamp := time.Now().Unix()
	if d.config.CertificateId == "" {
		// 创建证书
		createCertificateReq := &wangsucdn.CreateCertificateRequest{
			Timestamp:  timestamp,
			Name:       xtypes.ToPtr(fmt.Sprintf("certimate_%d", time.Now().UnixMilli())),
			AutoRenew:  xtypes.ToPtr("Off"),
			NewVersion: certificateNewVersionInfo,
		}
		createCertificateResp, err := d.sdkClient.CreateCertificate(createCertificateReq)
		d.logger.Debug("sdk request 'cdnpro.CreateCertificate'", slog.Any("request", createCertificateReq), slog.Any("response", createCertificateResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdnpro.CreateCertificate': %w", err)
		}

		wangsuCertUrl = createCertificateResp.CertificateLocation
		d.logger.Info("ssl certificate uploaded", slog.Any("certUrl", wangsuCertUrl))

		wangsuCertIdMatches := regexp.MustCompile(`/certificates/([a-zA-Z0-9-]+)`).FindStringSubmatch(wangsuCertUrl)
		if len(wangsuCertIdMatches) > 1 {
			wangsuCertId = wangsuCertIdMatches[1]
		}

		wangsuCertVer = 1
	} else {
		// 更新证书
		updateCertificateReq := &wangsucdn.UpdateCertificateRequest{
			Timestamp:  timestamp,
			Name:       xtypes.ToPtr(fmt.Sprintf("certimate_%d", time.Now().UnixMilli())),
			AutoRenew:  xtypes.ToPtr("Off"),
			NewVersion: certificateNewVersionInfo,
		}
		updateCertificateResp, err := d.sdkClient.UpdateCertificate(d.config.CertificateId, updateCertificateReq)
		d.logger.Debug("sdk request 'cdnpro.CreateCertificate'", slog.Any("certificateId", d.config.CertificateId), slog.Any("request", updateCertificateReq), slog.Any("response", updateCertificateResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdnpro.UpdateCertificate': %w", err)
		}

		wangsuCertUrl = updateCertificateResp.CertificateLocation
		d.logger.Info("ssl certificate uploaded", slog.Any("certUrl", wangsuCertUrl))

		wangsuCertIdMatches := regexp.MustCompile(`/certificates/([a-zA-Z0-9-]+)`).FindStringSubmatch(wangsuCertUrl)
		if len(wangsuCertIdMatches) > 1 {
			wangsuCertId = wangsuCertIdMatches[1]
		}

		wangsuCertVerMatches := regexp.MustCompile(`/versions/(\d+)`).FindStringSubmatch(wangsuCertUrl)
		if len(wangsuCertVerMatches) > 1 {
			n, _ := strconv.ParseInt(wangsuCertVerMatches[1], 10, 32)
			wangsuCertVer = int32(n)
		}
	}

	// 创建部署任务
	// REF: https://www.wangsu.com/document/api-doc/27034
	createDeploymentTaskReq := &wangsucdn.CreateDeploymentTaskRequest{
		Name:   xtypes.ToPtr(fmt.Sprintf("certimate_%d", time.Now().UnixMilli())),
		Target: xtypes.ToPtr(d.config.Environment),
		Actions: &[]wangsucdn.DeploymentTaskActionInfo{
			{
				Action:        xtypes.ToPtr("deploy_cert"),
				CertificateId: xtypes.ToPtr(wangsuCertId),
				Version:       xtypes.ToPtr(wangsuCertVer),
			},
		},
	}
	if d.config.WebhookId != "" {
		createDeploymentTaskReq.Webhook = xtypes.ToPtr(d.config.WebhookId)
	}
	createDeploymentTaskResp, err := d.sdkClient.CreateDeploymentTask(createDeploymentTaskReq)
	d.logger.Debug("sdk request 'cdnpro.CreateCertificate'", slog.Any("request", createDeploymentTaskReq), slog.Any("response", createDeploymentTaskResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdnpro.CreateDeploymentTask': %w", err)
	}

	// 循环获取部署任务详细信息，等待任务状态变更
	// REF: https://www.wangsu.com/document/api-doc/27038
	var wangsuTaskId string
	wangsuTaskMatches := regexp.MustCompile(`/deploymentTasks/([a-zA-Z0-9-]+)`).FindStringSubmatch(createDeploymentTaskResp.DeploymentTaskLocation)
	if len(wangsuTaskMatches) > 1 {
		wangsuTaskId = wangsuTaskMatches[1]
	}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		getDeploymentTaskDetailResp, err := d.sdkClient.GetDeploymentTaskDetail(wangsuTaskId)
		d.logger.Info("sdk request 'cdnpro.GetDeploymentTaskDetail'", slog.Any("taskId", wangsuTaskId), slog.Any("response", getDeploymentTaskDetailResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'cdnpro.GetDeploymentTaskDetail': %w", err)
		}

		if getDeploymentTaskDetailResp.Status == "failed" {
			return nil, errors.New("unexpected deployment task status")
		} else if getDeploymentTaskDetailResp.Status == "succeeded" || getDeploymentTaskDetailResp.FinishTime != "" {
			break
		}

		d.logger.Info(fmt.Sprintf("waiting for deployment task completion (current status: %s) ...", getDeploymentTaskDetailResp.Status))
		time.Sleep(time.Second * 5)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(accessKeyId, accessKeySecret string) (*wangsucdn.Client, error) {
	return wangsucdn.NewClient(accessKeyId, accessKeySecret)
}

func encryptPrivateKey(privkeyPEM string, apiKey string, timestamp int64) (string, error) {
	date := time.Unix(timestamp, 0).UTC()
	dateStr := date.Format("Mon, 02 Jan 2006 15:04:05 GMT")

	mac := hmac.New(sha256.New, []byte(apiKey))
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

	plainBytes := []byte(privkeyPEM)
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
