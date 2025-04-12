package dtos

type CertificateArchiveFileReq struct {
	CertificateId string `json:"-"`
	Format        string `json:"format"`
}

type CertificateArchiveFileResp struct {
	FileBytes  []byte `json:"fileBytes"`
	FileFormat string `json:"fileFormat"`
}

type CertificateValidateCertificateReq struct {
	Certificate string `json:"certificate"`
}

type CertificateValidateCertificateResp struct {
	IsValid bool   `json:"isValid"`
	Domains string `json:"domains,omitempty"`
}

type CertificateValidatePrivateKeyReq struct {
	PrivateKey string `json:"privateKey"`
}

type CertificateValidatePrivateKeyResp struct {
	IsValid bool `json:"isValid"`
}
