package dtos

import (
	"time"
)

type Environment struct {
	Type                     string    `json:"@type"`
	IdLink                   string    `json:"@id"`
	Id                       string    `json:"id"`
	PropertyID               string    `json:"property_id"`
	LegacyAccountNumber      string    `json:"legacy_account_number"`
	Name                     string    `json:"name"`
	CanMembersDeploy         bool      `json:"can_members_deploy"`
	OnlyMaintainersCanDeploy bool      `json:"only_maintainers_can_deploy"`
	HttpRequestLogging       bool      `json:"http_request_logging"`
	DefaultDomainName        string    `json:"default_domain_name"`
	PciCompliance            bool      `json:"pci_compliance"`
	DnsDomainName            string    `json:"dns_domain_name"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type EnvironmentsResponse struct {
	Type       string        `json:"@type"`
	Id         string        `json:"@id"`
	TotalItems int           `json:"total_items"`
	Items      []Environment `json:"items"`
}
