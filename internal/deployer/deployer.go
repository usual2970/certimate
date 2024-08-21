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

type DeployerOption struct {
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
	switch access.GetString("configType") {
	case configTypeAliyun:
		option := &DeployerOption{
			Domain:  record.GetString("domain"),
			Product: getProduct(record),
			Access:  access.GetString("config"),
			Certificate: applicant.Certificate{
				Certificate: record.GetString("certificate"),
				PrivateKey:  record.GetString("privateKey"),
			},
		}
		return NewAliyun(option)
	}
	return nil, errors.New("not implemented")
}

func getProduct(record *models.Record) string {
	targetType := record.GetString("targetType")
	rs := strings.Split(targetType, "-")
	return rs[1]
}
