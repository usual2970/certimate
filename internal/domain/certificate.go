package domain

import (
	"time"
	"os"
	"strconv"
)

// 获取环境变量"day"的值
dayStr := os.Getenv("day")

// 检查环境变量是否存在
if dayStr == "" {
	// 如果不存在，设置默认值为"10"
	dayStr = "10"
}

// 将字符串转换为整数
day, err := strconv.Atoi(dayStr)
if err != nil {
	// 如果转换失败，打印错误并设置默认值
	fmt.Println("Error converting day to integer:", err)
	day = 10
}

var ValidityDuration = time.Hour * 24 * day

type Certificate struct {
	Meta
	SAN               string    `json:"san" db:"san"`
	Certificate       string    `json:"certificate" db:"certificate"`
	PrivateKey        string    `json:"privateKey" db:"privateKey"`
	IssuerCertificate string    `json:"issuerCertificate" db:"issuerCertificate"`
	CertUrl           string    `json:"certUrl" db:"certUrl"`
	CertStableUrl     string    `json:"certStableUrl" db:"certStableUrl"`
	Output            string    `json:"output" db:"output"`
	Workflow          string    `json:"workflow" db:"workflow"`
	ExpireAt          time.Time `json:"ExpireAt" db:"expireAt"`
	NodeId            string    `json:"nodeId" db:"nodeId"`
}

type MetaData struct {
	Version            string              `json:"version"`
	SerialNumber       string              `json:"serialNumber"`
	Validity           CertificateValidity `json:"validity"`
	SignatureAlgorithm string              `json:"signatureAlgorithm"`
	Issuer             CertificateIssuer   `json:"issuer"`
	Subject            CertificateSubject  `json:"subject"`
}

type CertificateIssuer struct {
	Country      string `json:"country"`
	Organization string `json:"organization"`
	CommonName   string `json:"commonName"`
}

type CertificateSubject struct {
	CN string `json:"CN"`
}

type CertificateValidity struct {
	NotBefore string `json:"notBefore"`
	NotAfter  string `json:"notAfter"`
}
