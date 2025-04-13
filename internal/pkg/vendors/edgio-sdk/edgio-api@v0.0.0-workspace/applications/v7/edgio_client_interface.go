package edgio_api

import (
	"context"

	"github.com/Edgio/edgio-api/applications/v7/dtos"
)

type EdgioClientInterface interface {
	GetProperty(ctx context.Context, propertyID string) (*dtos.Property, error)
	GetProperties(page int, pageSize int, organizationID string) (*dtos.Properties, error)
	CreateProperty(ctx context.Context, organizationID, slug string) (*dtos.Property, error)
	DeleteProperty(propertyID string) error
	UpdateProperty(ctx context.Context, propertyID string, slug string) (*dtos.Property, error)
	GetEnvironments(page, pageSize int, propertyID string) (*dtos.EnvironmentsResponse, error)
	GetEnvironment(environmentID string) (*dtos.Environment, error)
	CreateEnvironment(propertyID, name string, onlyMaintainersCanDeploy, httpRequestLogging bool) (*dtos.Environment, error)
	UpdateEnvironment(environmentID, name string, onlyMaintainersCanDeploy, httpRequestLogging, preserveCache bool) (*dtos.Environment, error)
	DeleteEnvironment(environmentID string) error
	GetTlsCert(tlsCertId string) (*dtos.TLSCertResponse, error)
	UploadTlsCert(req dtos.UploadTlsCertRequest) (*dtos.TLSCertResponse, error)
	GenerateTlsCert(environmentId string) (*dtos.TLSCertResponse, error)
	GetTlsCerts(page int, pageSize int, environmentID string) (*dtos.TLSCertSResponse, error)
	UploadCdnConfiguration(config *dtos.CDNConfiguration) (*dtos.CDNConfiguration, error)
	GetCDNConfiguration(configID string) (*dtos.CDNConfiguration, error)
}
