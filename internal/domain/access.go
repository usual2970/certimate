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
