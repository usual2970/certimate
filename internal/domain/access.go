package domain

import (
	"encoding/json"
	"time"
)

const CollectionNameAccess = "access"

type Access struct {
	Meta
	Name      string     `json:"name" db:"name"`
	Provider  string     `json:"provider" db:"provider"`
	Config    string     `json:"config" db:"config"`
	Usage     string     `json:"usage" db:"usage"`
	DeletedAt *time.Time `json:"deleted" db:"deleted"`
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
	Mode     string `json:"mode,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
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

type AccessConfigForClouDNS struct {
	AuthId       string `json:"authId"`
	AuthPassword string `json:"authPassword"`
}

type AccessConfigForDogeCloud struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type AccessConfigForEdgio struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type AccessConfigForGname struct {
	AppId  string `json:"appId"`
	AppKey string `json:"appKey"`
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
	KubeConfig string `json:"kubeConfig,omitempty"`
}

type AccessConfigForNameDotCom struct {
	Username string `json:"username"`
	ApiToken string `json:"apiToken"`
}

type AccessConfigForNameSilo struct {
	ApiKey string `json:"apiKey"`
}

type AccessConfigForNS1 struct {
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

type AccessConfigForRainYun struct {
	ApiKey string `json:"apiKey"`
}

type AccessConfigForSSH struct {
	Host          string `json:"host"`
	Port          int32  `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password,omitempty"`
	Key           string `json:"key,omitempty"`
	KeyPassphrase string `json:"keyPassphrase,omitempty"`
}

type AccessConfigForTencentCloud struct {
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
}

type AccessConfigForUCloud struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	ProjectId  string `json:"projectId,omitempty"`
}

type AccessConfigForVolcEngine struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type AccessConfigForWebhook struct {
	Url string `json:"url"`
}

type AccessConfigForWestcn struct {
	Username    string `json:"username"`
	ApiPassword string `json:"password"`
}
