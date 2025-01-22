package dtos

type CertificateArchiveFileReq struct {
	CertificateId string `json:"-"`
	Format        string `json:"format"`
}

type CertificateArchiveFileResp struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
}

type CertificateValidateCertificateReq struct {
	Certificate string `json:"certificate"`
}

type CertificateValidateCertificateResp struct {
	Domains string `json:"domains"`
}

type CertificateValidatePrivateKeyReq struct {
	PrivateKey string `json:"privateKey"`
}

type CertificateUploadReq struct {
	WorkflowId     string `json:"workflowId"`
	WorkflowNodeId string `json:"workflowNodeId"`
	CertificateId  string `json:"certificateId"`
	Certificate    string `json:"certificate"`
	PrivateKey     string `json:"privateKey"`
}
