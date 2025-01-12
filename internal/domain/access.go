package domain

import (
	"encoding/json"
	"time"
)

type Access struct {
	Meta
	Name      string    `json:"name" db:"name"`
	Provider  string    `json:"provider" db:"provider"`
	Config    string    `json:"config" db:"config"`
	Usage     string    `json:"usage" db:"usage"`
	DeletedAt time.Time `json:"deleted" db:"deleted"`
}

func (a *Access) UnmarshalConfigToMap() (map[string]any, error) {
	config := make(map[string]any)
	if err := json.Unmarshal([]byte(a.Config), &config); err != nil {
		return nil, err
	}

	return config, nil
}

type AccessConfigForACMEHttpReq struct {
	Endpoint string `json:"endpoint"`
	Mode     string `json:"mode"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccessConfigForAliyun struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type AccessConfigForAWS struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type AccessConfigForAzure struct {
	TenantId     string `json:"tenantId"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	CloudName    string `json:"cloudName,omitempty"`
}

type AccessConfigForBaiduCloud struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type AccessConfigForBytePlus struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type AccessConfigForCloudflare struct {
	DnsApiToken string `json:"dnsApiToken"`
}

type AccessConfigForDogeCloud struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type AccessConfigForGoDaddy struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}

type AccessConfigForHuaweiCloud struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type AccessConfigForLocal struct{}

type AccessConfigForKubernetes struct {
	KubeConfig string `json:"kubeConfig"`
}

type AccessConfigForNameDotCom struct {
	Username string `json:"username"`
	ApiToken string `json:"apiToken"`
}

type AccessConfigForNameSilo struct {
	ApiKey string `json:"apiKey"`
}

type AccessConfigForPowerDNS struct {
	ApiUrl string `json:"apiUrl"`
	ApiKey string `json:"apiKey"`
}

type AccessConfigForQiniu struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type AccessConfigForSSH struct {
	Host          string `json:"host"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Key           string `json:"key"`
	KeyPassphrase string `json:"keyPassphrase"`
}

type AccessConfigForTencentCloud struct {
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
}

type AccessConfigForVolcEngine struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type AccessConfigForWebhook struct {
	Url string `json:"url"`
}
