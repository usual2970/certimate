package cert

import (
	"github.com/baidubce/bce-sdk-go/services/cert"
)

type Client struct {
	*cert.Client
}

func NewClient(ak, sk, endPoint string) (*Client, error) {
	client, err := cert.NewClient(ak, sk, endPoint)
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}
