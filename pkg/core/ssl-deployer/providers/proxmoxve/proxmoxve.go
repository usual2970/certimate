package proxmoxve

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/luthermonson/go-proxmox"

	"github.com/certimate-go/certimate/pkg/core"
	xhttp "github.com/certimate-go/certimate/pkg/utils/http"
)

type SSLDeployerProviderConfig struct {
	// Proxmox VE 服务地址。
	ServerUrl string `json:"serverUrl"`
	// Proxmox VE API Token。
	ApiToken string `json:"apiToken"`
	// Proxmox VE API Token Secret。
	ApiTokenSecret string `json:"apiTokenSecret,omitempty"`
	// 是否允许不安全的连接。
	AllowInsecureConnections bool `json:"allowInsecureConnections,omitempty"`
	// 集群节点名称。
	NodeName string `json:"nodeName"`
	// 是否自动重启。
	AutoRestart bool `json:"autoRestart"`
}

type SSLDeployerProvider struct {
	config    *SSLDeployerProviderConfig
	logger    *slog.Logger
	sdkClient *proxmox.Client
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	client, err := createSDKClient(config.ServerUrl, config.ApiToken, config.ApiTokenSecret, config.AllowInsecureConnections)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	return &SSLDeployerProvider{
		config:    config,
		logger:    slog.Default(),
		sdkClient: client,
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	if d.config.NodeName == "" {
		return nil, errors.New("config `nodeName` is required")
	}

	// 获取节点信息
	// REF: https://pve.proxmox.com/pve-docs/api-viewer/index.html#/nodes/{node}
	node, err := d.sdkClient.Node(context.TODO(), d.config.NodeName)
	if err != nil {
		return nil, fmt.Errorf("failed to get node '%s': %w", d.config.NodeName, err)
	}

	// 上传自定义证书
	// REF: https://pve.proxmox.com/pve-docs/api-viewer/index.html#/nodes/{node}/certificates/custom
	err = node.UploadCustomCertificate(context.TODO(), &proxmox.CustomCertificate{
		Certificates: certPEM,
		Key:          privkeyPEM,
		Force:        true,
		Restart:      d.config.AutoRestart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload custom certificate to node '%s': %w", node.Name, err)
	}

	return &core.SSLDeployResult{}, nil
}

func createSDKClient(serverUrl, apiToken, apiTokenSecret string, skipTlsVerify bool) (*proxmox.Client, error) {
	if _, err := url.Parse(serverUrl); err != nil {
		return nil, errors.New("invalid pve server url")
	}

	if apiToken == "" {
		return nil, errors.New("invalid pve api token")
	}

	httpClient := &http.Client{
		Transport: xhttp.NewDefaultTransport(),
		Timeout:   http.DefaultClient.Timeout,
	}
	if skipTlsVerify {
		transport := xhttp.NewDefaultTransport()
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		transport.TLSClientConfig.InsecureSkipVerify = true
		httpClient.Transport = transport
	}
	client := proxmox.NewClient(
		strings.TrimRight(serverUrl, "/")+"/api2/json",
		proxmox.WithHTTPClient(httpClient),
		proxmox.WithAPIToken(apiToken, apiTokenSecret),
	)

	return client, nil
}
