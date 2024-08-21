package deployer

import (
	"certimate/internal/applicant"
	"context"
	"errors"
	"strings"

	"github.com/pocketbase/pocketbase/models"
)

const (
	configTypeAliyun = "aliyun"
)

const (
	targetAliyunOss = "aliyun-oss"
	targetAliyunCdn = "aliyun-cdn"
)

type DeployerOption struct {
	DomainId    string                `json:"domainId"`
	Domain      string                `json:"domain"`
	Product     string                `json:"product"`
	Access      string                `json:"access"`
	Certificate applicant.Certificate `json:"certificate"`
}

type Deployer interface {
	Deploy(ctx context.Context) error
}

func Get(record *models.Record) (Deployer, error) {
	access := record.ExpandedOne("targetAccess")
	option := &DeployerOption{
		DomainId: record.Id,
		Domain:   record.GetString("domain"),
		Product:  getProduct(record),
		Access:   access.GetString("config"),
		Certificate: applicant.Certificate{
			Certificate: record.GetString("certificate"),
			PrivateKey:  record.GetString("privateKey"),
		},
	}
	switch record.GetString("targetType") {
	case targetAliyunOss:
		return NewAliyun(option)
	case targetAliyunCdn:
		return NewAliyunCdn(option)
	}
	return nil, errors.New("not implemented")
}

func getProduct(record *models.Record) string {
	targetType := record.GetString("targetType")
	rs := strings.Split(targetType, "-")
	return rs[1]
}
