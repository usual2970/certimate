package domain

import "time"

type Access struct {
	Meta
	Name      string    `json:"name" db:"name"`
	Provider  string    `json:"provider" db:"provider"`
	Config    string    `json:"config" db:"config"`
	Usage     string    `json:"usage" db:"usage"`
	DeletedAt time.Time `json:"deleted" db:"deleted"`
}

type AccessProviderType string

/*
提供商类型常量值。

	注意：如果追加新的常量值，请保持以 ASCII 排序。
	NOTICE: If you add new constant, please keep ASCII order.
*/
const (
	ACCESS_PROVIDER_ACMEHTTPREQ  = AccessProviderType("acmehttpreq")
	ACCESS_PROVIDER_ALIYUN       = AccessProviderType("aliyun")
	ACCESS_PROVIDER_AWS          = AccessProviderType("aws")
	ACCESS_PROVIDER_BAIDUCLOUD   = AccessProviderType("baiducloud")
	ACCESS_PROVIDER_BYTEPLUS     = AccessProviderType("byteplus")
	ACCESS_PROVIDER_CLOUDFLARE   = AccessProviderType("cloudflare")
	ACCESS_PROVIDER_DOGECLOUD    = AccessProviderType("dogecloud")
	ACCESS_PROVIDER_GODADDY      = AccessProviderType("godaddy")
	ACCESS_PROVIDER_HUAWEICLOUD  = AccessProviderType("huaweicloud")
	ACCESS_PROVIDER_KUBERNETES   = AccessProviderType("k8s")
	ACCESS_PROVIDER_LOCAL        = AccessProviderType("local")
	ACCESS_PROVIDER_NAMEDOTCOM   = AccessProviderType("namedotcom")
	ACCESS_PROVIDER_NAMESILO     = AccessProviderType("namesilo")
	ACCESS_PROVIDER_POWERDNS     = AccessProviderType("powerdns")
	ACCESS_PROVIDER_QINIU        = AccessProviderType("qiniu")
	ACCESS_PROVIDER_SSH          = AccessProviderType("ssh")
	ACCESS_PROVIDER_TENCENTCLOUD = AccessProviderType("tencentcloud")
	ACCESS_PROVIDER_VOLCENGINE   = AccessProviderType("volcengine")
	ACCESS_PROVIDER_WEBHOOK      = AccessProviderType("webhook")
)

type ACMEHttpReqAccessConfig struct {
	Endpoint string `json:"endpoint"`
	Mode     string `json:"mode"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AliyunAccessConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type AWSAccessConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Region          string `json:"region"`
	HostedZoneId    string `json:"hostedZoneId"`
}

type BaiduCloudAccessConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type BytePlusAccessConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type CloudflareAccessConfig struct {
	DnsApiToken string `json:"dnsApiToken"`
}

type DogeCloudAccessConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type GoDaddyAccessConfig struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}

type HuaweiCloudAccessConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Region          string `json:"region"`
}

type LocalAccessConfig struct{}

type KubernetesAccessConfig struct {
	KubeConfig string `json:"kubeConfig"`
}

type NameDotComAccessConfig struct {
	Username string `json:"username"`
	ApiToken string `json:"apiToken"`
}

type NameSiloAccessConfig struct {
	ApiKey string `json:"apiKey"`
}

type PowerDNSAccessConfig struct {
	ApiUrl string `json:"apiUrl"`
	ApiKey string `json:"apiKey"`
}

type QiniuAccessConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type SSHAccessConfig struct {
	Host          string `json:"host"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Key           string `json:"key"`
	KeyPassphrase string `json:"keyPassphrase"`
}

type TencentCloudAccessConfig struct {
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
}

type VolcEngineAccessConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type WebhookAccessConfig struct {
	Url string `json:"url"`
}
