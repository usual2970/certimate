package domain

type AliyunAccess struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type TencentAccess struct {
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
}

type HuaweiCloudAccess struct {
	Region          string `json:"region"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type AwsAccess struct {
	Region          string `json:"region"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	HostedZoneId    string `json:"hostedZoneId"`
}

type CloudflareAccess struct {
	DnsApiToken string `json:"dnsApiToken"`
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

type PdnsAccess struct {
	ApiUrl string `json:"apiUrl"`
	ApiKey string `json:"apiKey"`
}

type HttpreqAccess struct {
	Endpoint string `json:"endpoint"`
	Mode     string `json:"mode"`
	Username string `json:"username"`
	Password string `json:"password"`
}
