package lego_baiducloud

import (
	"errors"
	"fmt"
	"strings"
	"time"

	bceDns "github.com/baidubce/bce-sdk-go/services/dns"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"
	"github.com/google/uuid"
)

const (
	envNamespace = "BAIDUCLOUD_"

	EnvAccessKeyID     = envNamespace + "ACCESS_KEY_ID"
	EnvSecretAccessKey = envNamespace + "SECRET_ACCESS_KEY"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvHTTPTimeout        = envNamespace + "HTTP_TIMEOUT"
)

var _ challenge.ProviderTimeout = (*DNSProvider)(nil)

type Config struct {
	AccessKeyID     string
	SecretAccessKey string

	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int32
	HTTPTimeout        time.Duration
}

type DNSProvider struct {
	client *bceDns.Client
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
	values, err := env.Get(EnvAccessKeyID, EnvSecretAccessKey)
	if err != nil {
		return nil, fmt.Errorf("baiducloud: %w", err)
	}

	config := NewDefaultConfig()
	config.AccessKeyID = values[EnvAccessKeyID]
	config.SecretAccessKey = values[EnvSecretAccessKey]

	return NewDNSProviderConfig(config)
}

func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("baiducloud: the configuration of the DNS provider is nil")
	}

	client, err := bceDns.NewClient(config.AccessKeyID, config.SecretAccessKey, "")
	if err != nil {
		return nil, err
	} else {
		if client.Config != nil {
			client.Config.ConnectionTimeoutInMillis = int(config.HTTPTimeout.Milliseconds())
		}
	}

	return &DNSProvider{
		client: client,
		config: config,
	}, nil
}

func (d *DNSProvider) Present(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	zoneName, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("baiducloud: %w", err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, zoneName)
	if err != nil {
		return fmt.Errorf("baiducloud: %w", err)
	}

	if err := d.addOrUpdateDNSRecord(domain, subDomain, info.Value); err != nil {
		return fmt.Errorf("baiducloud: %w", err)
	}

	return nil
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	subDomain := dns01.UnFqdn(fqdn)

	if err := d.removeDNSRecord(domain, subDomain, value); err != nil {
		return fmt.Errorf("baiducloud: %w", err)
	}

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}

func (d *DNSProvider) getDNSRecord(domain, subDomain string) (*bceDns.Record, error) {
	pageMarker := ""
	pageSize := 1000
	for {
		request := &bceDns.ListRecordRequest{}
		request.Rr = domain
		request.Marker = pageMarker
		request.MaxKeys = pageSize

		response, err := d.client.ListRecord(domain, request)
		if err != nil {
			return nil, err
		}

		for _, record := range response.Records {
			if record.Type == "TXT" && record.Rr == subDomain {
				return &record, nil
			}
		}

		if len(response.Records) < pageSize {
			break
		}

		pageMarker = response.NextMarker
	}

	return nil, nil
}

func (d *DNSProvider) addOrUpdateDNSRecord(domain, subDomain, value string) error {
	record, err := d.getDNSRecord(domain, subDomain)
	if err != nil {
		return err
	}

	if record == nil {
		request := &bceDns.CreateRecordRequest{
			Type:  "TXT",
			Rr:    subDomain,
			Value: value,
			Ttl:   &d.config.TTL,
		}
		err := d.client.CreateRecord(domain, request, d.generateClientToken())
		return err
	} else {
		request := &bceDns.UpdateRecordRequest{
			Type:  "TXT",
			Rr:    subDomain,
			Value: value,
			Ttl:   &d.config.TTL,
		}
		err := d.client.UpdateRecord(domain, record.Id, request, d.generateClientToken())
		return err
	}
}

func (d *DNSProvider) removeDNSRecord(domain, subDomain, value string) error {
	record, err := d.getDNSRecord(domain, subDomain)
	if err != nil {
		return err
	}

	if record == nil {
		return nil
	}

	err = d.client.DeleteRecord(domain, record.Id, d.generateClientToken())
	return err
}

func (d *DNSProvider) generateClientToken() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
