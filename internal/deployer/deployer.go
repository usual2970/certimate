package deployer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase/models"

	"github.com/usual2970/certimate/internal/applicant"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/app"
)

const (
	targetAliyunOSS  = "aliyun-oss"
	targetAliyunCDN  = "aliyun-cdn"
	targetAliyunESA  = "aliyun-dcdn"
	targetTencentCDN = "tencent-cdn"
	targetQiniuCdn   = "qiniu-cdn"
	targetLocal      = "local"
	targetSSH        = "ssh"
	targetWebhook    = "webhook"
	targetK8sSecret  = "k8s-secret"
)

type DeployerOption struct {
	DomainId     string                `json:"domainId"`
	Domain       string                `json:"domain"`
	Product      string                `json:"product"`
	Access       string                `json:"access"`
	AceessRecord *models.Record        `json:"-"`
	DeployConfig domain.DeployConfig   `json:"deployConfig"`
	Certificate  applicant.Certificate `json:"certificate"`
	Variables    map[string]string     `json:"variables"`
}

type Deployer interface {
	Deploy(ctx context.Context) error
	GetInfo() []string
	GetID() string
}

func Gets(record *models.Record, cert *applicant.Certificate) ([]Deployer, error) {
	rs := make([]Deployer, 0)
	if record.GetString("deployConfig") == "" {
		return rs, nil
	}

	deployConfigs := make([]domain.DeployConfig, 0)

	err := record.UnmarshalJSONField("deployConfig", &deployConfigs)
	if err != nil {
		return nil, fmt.Errorf("解析部署配置失败: %w", err)
	}

	if len(deployConfigs) == 0 {
		return rs, nil
	}

	for _, deployConfig := range deployConfigs {
		deployer, err := getWithDeployConfig(record, cert, deployConfig)
		if err != nil {
			return nil, err
		}

		rs = append(rs, deployer)
	}

	return rs, nil
}

func getWithDeployConfig(record *models.Record, cert *applicant.Certificate, deployConfig domain.DeployConfig) (Deployer, error) {
	access, err := app.GetApp().Dao().FindRecordById("access", deployConfig.Access)
	if err != nil {
		return nil, fmt.Errorf("access record not found: %w", err)
	}

	option := &DeployerOption{
		DomainId:     record.Id,
		Domain:       record.GetString("domain"),
		Product:      getProduct(deployConfig.Type),
		Access:       access.GetString("config"),
		AceessRecord: access,
		DeployConfig: deployConfig,
	}
	if cert != nil {
		option.Certificate = *cert
	} else {
		option.Certificate = applicant.Certificate{
			Certificate: record.GetString("certificate"),
			PrivateKey:  record.GetString("privateKey"),
		}
	}

	switch deployConfig.Type {
	case targetAliyunOSS:
		return NewAliyunOssDeployer(option)
	case targetAliyunCDN:
		return NewAliyunCdnDeployer(option)
	case targetAliyunESA:
		return NewAliyunEsaDeployer(option)
	case targetTencentCDN:
		return NewTencentCDNDeployer(option)
	case targetQiniuCdn:
		return NewQiniuCDNDeployer(option)
	case targetLocal:
		return NewLocalDeployer(option)
	case targetSSH:
		return NewSSHDeployer(option)
	case targetWebhook:
		return NewWebhookDeployer(option)
	case targetK8sSecret:
		return NewK8sSecretDeployer(option)
	}
	return nil, errors.New("not implemented")
}

func getProduct(t string) string {
	rs := strings.Split(t, "-")
	if len(rs) < 2 {
		return ""
	}
	return rs[1]
}

func toStr(tag string, data any) string {
	if data == nil {
		return tag
	}
	byts, _ := json.Marshal(data)
	return tag + "：" + string(byts)
}

func getDeployString(conf domain.DeployConfig, key string) string {
	if _, ok := conf.Config[key]; !ok {
		return ""
	}

	val, ok := conf.Config[key].(string)
	if !ok {
		return ""
	}

	return val
}

func getDeployVariables(conf domain.DeployConfig) map[string]string {
	rs := make(map[string]string)
	data, ok := conf.Config["variables"]
	if !ok {
		return rs
	}

	bts, _ := json.Marshal(data)

	kvData := make([]domain.KV, 0)

	if err := json.Unmarshal(bts, &kvData); err != nil {
		return rs
	}

	for _, kv := range kvData {
		rs[kv.Key] = kv.Value
	}

	return rs
}
