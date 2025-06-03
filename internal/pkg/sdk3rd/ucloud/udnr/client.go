package udnr

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

type UDNRClient struct {
	*ucloud.Client
}

func NewClient(config *ucloud.Config, credential *auth.Credential) *UDNRClient {
	meta := ucloud.ClientMeta{Product: "UDNR"}
	client := ucloud.NewClientWithMeta(config, credential, meta)
	return &UDNRClient{
		client,
	}
}
