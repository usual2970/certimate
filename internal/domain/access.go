package domain

type AliyunAccess struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type TencentAccess struct {
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
}

type CloudflareAccess struct {
	DnsApiToken string `json:"dnsApiToken"`
}

type AwsAccess struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SSLprovider     string `json:"sslprovider"`
}

type QiniuAccess struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type NameSiloAccess struct {
	ApiKey string `json:"apiKey"`
}

type GodaddyAccess struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}
