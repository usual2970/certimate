package deployer

import (
	"certimate/internal/applicant"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/pocketbase/pocketbase/models"
)

const (
	configTypeAliyun = "aliyun"
)

const (
	targetAliyunOss  = "aliyun-oss"
	targetAliyunCdn  = "aliyun-cdn"
	targetSSH        = "ssh"
	targetWebhook    = "webhook"
	targetTencentCdn = "tencent-cdn"
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
	GetInfo() []string
}

func Get(record *models.Record, cert *applicant.Certificate) (Deployer, error) {
	access := record.ExpandedOne("targetAccess")
	option := &DeployerOption{
		DomainId: record.Id,
		Domain:   record.GetString("domain"),
		Product:  getProduct(record),
		Access:   access.GetString("config"),
	}
	if cert != nil {
		option.Certificate = *cert
	} else {
		option.Certificate = applicant.Certificate{
			Certificate: record.GetString("certificate"),
			PrivateKey:  record.GetString("privateKey"),
		}
	}

	switch record.GetString("targetType") {
	case targetAliyunOss:
		return NewAliyun(option)
	case targetAliyunCdn:
		return NewAliyunCdn(option)
	case targetSSH:
		return NewSSH(option)
	case targetWebhook:
		return NewWebhook(option)
	case targetTencentCdn:
		return NewTencentCdn(option)
	}
	return nil, errors.New("not implemented")
}

func getProduct(record *models.Record) string {
	targetType := record.GetString("targetType")
	rs := strings.Split(targetType, "-")
	if len(rs) < 2 {
		return ""
	}
	return rs[1]
}

func toStr(tag string, data any) string {
	if data == nil {
		return tag
	}
	byts, _ := json.Marshal(data)
	return tag + "：" + string(byts)
}
