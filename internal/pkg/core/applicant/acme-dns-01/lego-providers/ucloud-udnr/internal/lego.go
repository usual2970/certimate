package internal

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"

	"github.com/usual2970/certimate/internal/pkg/sdk3rd/ucloud/udnr"
)

const (
	envNamespace = "UCLOUDUDNR_"

	EnvPublicKey  = envNamespace + "PUBLIC_KEY"
	EnvPrivateKey = envNamespace + "PRIVATE_KEY"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvHTTPTimeout        = envNamespace + "HTTP_TIMEOUT"
)

var _ challenge.ProviderTimeout = (*DNSProvider)(nil)

type Config struct {
	PrivateKey string
	PublicKey  string

	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int32
	HTTPTimeout        time.Duration
}

type DNSProvider struct {
	client *udnr.UDNRClient
	config *Config
}

func NewDefaultConfig() *Config {
	return &Config{
		TTL:                int32(env.GetOrDefaultInt(EnvTTL, 300)),
		PropagationTimeout: env.GetOrDefaultSecond(EnvPropagationTimeout, 2*time.Minute),
		PollingInterval:    env.GetOrDefaultSecond(EnvPollingInterval, dns01.DefaultPollingInterval),
		HTTPTimeout:        env.GetOrDefaultSecond(EnvHTTPTimeout, 30*time.Second),
	}
}

func NewDNSProvider() (*DNSProvider, error) {
	values, err := env.Get(EnvPrivateKey, EnvPublicKey)
	if err != nil {
		return nil, fmt.Errorf("ucloud-udnr: %w", err)
	}

	config := NewDefaultConfig()
	config.PrivateKey = values[EnvPrivateKey]
	config.PublicKey = values[EnvPublicKey]

	return NewDNSProviderConfig(config)
}

func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("ucloud-udnr: the configuration of the DNS provider is nil")
	}

	cfg := ucloud.NewConfig()
	credential := auth.NewCredential()
	credential.PrivateKey = config.PrivateKey
	credential.PublicKey = config.PublicKey
	client := udnr.NewClient(&cfg, &credential)

	return &DNSProvider{
		client: client,
		config: config,
	}, nil
}

func (d *DNSProvider) Present(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("ucloud-udnr: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("ucloud-udnr: %w", err)
	}

	udnrDomainDNSQueryReq := d.client.NewQueryDomainDNSRequest()
	udnrDomainDNSQueryReq.Dn = ucloud.String(authZone)
	if udnrDomainDNSQueryResp, err := d.client.QueryDomainDNS(udnrDomainDNSQueryReq); err != nil {
		return fmt.Errorf("ucloud-udnr: %w", err)
	} else {
		for _, record := range udnrDomainDNSQueryResp.Data {
			if record.DnsType == "TXT" && record.RecordName == subDomain {
				udnrDomainDNSDeleteReq := d.client.NewDeleteDomainDNSRequest()
				udnrDomainDNSDeleteReq.Dn = ucloud.String(authZone)
				udnrDomainDNSDeleteReq.DnsType = ucloud.String(record.DnsType)
				udnrDomainDNSDeleteReq.RecordName = ucloud.String(record.RecordName)
				udnrDomainDNSDeleteReq.Content = ucloud.String(record.Content)
				d.client.DeleteDomainDNS(udnrDomainDNSDeleteReq)
				break
			}
		}
	}

	udnrDomainDNSAddReq := d.client.NewAddDomainDNSRequest()
	udnrDomainDNSAddReq.Dn = ucloud.String(authZone)
	udnrDomainDNSAddReq.DnsType = ucloud.String("TXT")
	udnrDomainDNSAddReq.RecordName = ucloud.String(subDomain)
	udnrDomainDNSAddReq.Content = ucloud.String(info.Value)
	udnrDomainDNSAddReq.TTL = ucloud.Int(int(d.config.TTL))
	if _, err := d.client.AddDomainDNS(udnrDomainDNSAddReq); err != nil {
		return fmt.Errorf("ucloud-udnr: %w", err)
	}

	return nil
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("ucloud-udnr: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("ucloud-udnr: %w", err)
	}

	udnrDomainDNSQueryReq := d.client.NewQueryDomainDNSRequest()
	udnrDomainDNSQueryReq.Dn = ucloud.String(authZone)
	if udnrDomainDNSQueryResp, err := d.client.QueryDomainDNS(udnrDomainDNSQueryReq); err != nil {
		return fmt.Errorf("ucloud-udnr: %w", err)
	} else {
		for _, record := range udnrDomainDNSQueryResp.Data {
			if record.DnsType == "TXT" && record.RecordName == subDomain {
				udnrDomainDNSDeleteReq := d.client.NewDeleteDomainDNSRequest()
				udnrDomainDNSDeleteReq.Dn = ucloud.String(authZone)
				udnrDomainDNSDeleteReq.DnsType = ucloud.String(record.DnsType)
				udnrDomainDNSDeleteReq.RecordName = ucloud.String(record.RecordName)
				udnrDomainDNSDeleteReq.Content = ucloud.String(record.Content)
				d.client.DeleteDomainDNS(udnrDomainDNSDeleteReq)
				break
			}
		}
	}

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}
