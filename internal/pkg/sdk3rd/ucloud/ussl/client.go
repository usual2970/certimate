package ussl

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

type USSLClient struct {
	*ucloud.Client
}

func NewClient(config *ucloud.Config, credential *auth.Credential) *USSLClient {
	meta := ucloud.ClientMeta{Product: "USSL"}
	client := ucloud.NewClientWithMeta(config, credential, meta)
	return &USSLClient{
		client,
	}
}
