package deployer

import (
	"certimate/internal/applicant"
	"certimate/internal/utils/app"
	"certimate/internal/utils/variables"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase/models"
)

const (
	targetAliyunOss  = "aliyun-oss"
	targetAliyunCdn  = "aliyun-cdn"
	targetAliyunEsa  = "aliyun-dcdn"
	targetSSH        = "ssh"
	targetWebhook    = "webhook"
	targetTencentCdn = "tencent-cdn"
	targetQiniuCdn   = "qiniu-cdn"
	targetLocal      = "local"
)

type DeployerOption struct {
	DomainId     string                `json:"domainId"`
	Domain       string                `json:"domain"`
	Product      string                `json:"product"`
	Access       string                `json:"access"`
	AceessRecord *models.Record        `json:"-"`
	Certificate  applicant.Certificate `json:"certificate"`
	Variables    map[string]string     `json:"variables"`
}

type Deployer interface {
	Deploy(ctx context.Context) error
	GetInfo() []string
	GetID() string
}

func Gets(record *models.Record, cert *applicant.Certificate) ([]Deployer, error) {
	rs := make([]Deployer, 0)

	if record.GetString("targetAccess") != "" {
		singleDeployer, err := Get(record, cert)
		if err != nil {
			return nil, err
		}
		rs = append(rs, singleDeployer)
	}

	if record.GetString("group") != "" {
		group := record.ExpandedOne("group")

		if errs := app.GetApp().Dao().ExpandRecord(group, []string{"access"}, nil); len(errs) > 0 {

			errList := make([]error, 0)
			for name, err := range errs {
				errList = append(errList, fmt.Errorf("展开记录失败,%s: %w", name, err))
			}
			err := errors.Join(errList...)
			return nil, err
		}

		records := group.ExpandedAll("access")

		deployers, err := getByGroup(record, cert, records...)
		if err != nil {
			return nil, err
		}

		rs = append(rs, deployers...)

	}

	return rs, nil

}

func getByGroup(record *models.Record, cert *applicant.Certificate, accesses ...*models.Record) ([]Deployer, error) {

	rs := make([]Deployer, 0)

	for _, access := range accesses {
		deployer, err := getWithAccess(record, cert, access)
		if err != nil {
			return nil, err
		}
		rs = append(rs, deployer)
	}

	return rs, nil

}

func getWithAccess(record *models.Record, cert *applicant.Certificate, access *models.Record) (Deployer, error) {

	option := &DeployerOption{
		DomainId:     record.Id,
		Domain:       record.GetString("domain"),
		Product:      getProduct(record),
		Access:       access.GetString("config"),
		AceessRecord: access,
		Variables:    variables.Parse2Map(record.GetString("variables")),
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
	case targetAliyunEsa:
		return NewAliyunEsa(option)
	case targetSSH:
		return NewSSH(option)
	case targetWebhook:
		return NewWebhook(option)
	case targetTencentCdn:
		return NewTencentCdn(option)
	case targetQiniuCdn:

		return NewQiNiu(option)
	case targetLocal:
		return NewLocal(option), nil
	}
	return nil, errors.New("not implemented")
}

func Get(record *models.Record, cert *applicant.Certificate) (Deployer, error) {

	access := record.ExpandedOne("targetAccess")

	return getWithAccess(record, cert, access)
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
