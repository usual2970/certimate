package domain

import (
	"time"
)

const CollectionNameAccess = "access"

type Access struct {
	Meta
	Name      string         `json:"name" db:"name"`
	Provider  string         `json:"provider" db:"provider"`
	Config    map[string]any `json:"config" db:"config"`
	DeletedAt *time.Time     `json:"deleted" db:"deleted"`
}

type AccessConfigFor1Panel struct {
	ApiUrl                   string `json:"apiUrl"`
	ApiKey                   string `json:"apiKey"`
	AllowInsecureConnections bool   `json:"allowInsecureConnections,omitempty"`
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

type AccessConfigForBaishan struct {
	ApiToken string `json:"apiToken"`
}

type AccessConfigForBaotaPanel struct {
	ApiUrl                   string `json:"apiUrl"`
	ApiKey                   string `json:"apiKey"`
	AllowInsecureConnections bool   `json:"allowInsecureConnections,omitempty"`
}

type AccessConfigForBytePlus struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type AccessConfigForCacheFly struct {
	ApiToken string `json:"apiToken"`
}

type AccessConfigForCdnfly struct {
	ApiUrl    string `json:"apiUrl"`
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}

type AccessConfigForCloudflare struct {
	DnsApiToken string `json:"dnsApiToken"`
}

type AccessConfigForClouDNS struct {
	AuthId       string `json:"authId"`
	AuthPassword string `json:"authPassword"`
}

type AccessConfigForCMCCCloud struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type AccessConfigForDeSEC struct {
	Token string `json:"token"`
}

type AccessConfigForDNSLA struct {
	ApiId     string `json:"apiId"`
	ApiSecret string `json:"apiSecret"`
}

type AccessConfigForDogeCloud struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type AccessConfigForDynv6 struct {
	HttpToken string `json:"httpToken"`
}

type AccessConfigForEdgio struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type AccessConfigForGcore struct {
	ApiToken string `json:"apiToken"`
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

type AccessConfigForJDCloud struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type AccessConfigForKubernetes struct {
	KubeConfig string `json:"kubeConfig,omitempty"`
}

type AccessConfigForLocal struct{}

type AccessConfigForNamecheap struct {
	Username string `json:"username"`
	ApiKey   string `json:"apiKey"`
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

type AccessConfigForPorkbun struct {
	ApiKey       string `json:"apiKey"`
	SecretApiKey string `json:"secretApiKey"`
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

type AccessConfigForSafeLine struct {
	ApiUrl                   string `json:"apiUrl"`
	ApiToken                 string `json:"apiToken"`
	AllowInsecureConnections bool   `json:"allowInsecureConnections,omitempty"`
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

type AccessConfigForUpyun struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccessConfigForVercel struct {
	ApiAccessToken string `json:"apiAccessToken"`
	TeamId         string `json:"teamId,omitempty"`
}

type AccessConfigForVolcEngine struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type AccessConfigForWebhook struct {
	Url                      string `json:"url"`
	AllowInsecureConnections bool   `json:"allowInsecureConnections,omitempty"`
}

type AccessConfigForWestcn struct {
	Username    string `json:"username"`
	ApiPassword string `json:"password"`
}
