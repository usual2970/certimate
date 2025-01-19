package dtos

type CertificateArchiveFileReq struct {
	CertificateId string `json:"-"`
	Format        string `json:"format"`
}
