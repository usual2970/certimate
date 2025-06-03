package internal

import (
	"errors"
	"fmt"
	"strings"
	"time"

	bcedns "github.com/baidubce/bce-sdk-go/services/dns"
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
	client *bcedns.Client
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

	client, err := bcedns.NewClient(config.AccessKeyID, config.SecretAccessKey, "")
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

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("baiducloud: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("baiducloud: %w", err)
	}

	if err := d.addOrUpdateDNSRecord(dns01.UnFqdn(authZone), subDomain, info.Value); err != nil {
		return fmt.Errorf("baiducloud: %w", err)
	}

	return nil
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("baiducloud: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("baiducloud: %w", err)
	}

	if err := d.removeDNSRecord(dns01.UnFqdn(authZone), subDomain); err != nil {
		return fmt.Errorf("baiducloud: %w", err)
	}

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}

func (d *DNSProvider) findDNSRecord(zoneName, subDomain string) (*bcedns.Record, error) {
	pageMarker := ""
	pageSize := 1000
	for {
		request := &bcedns.ListRecordRequest{}
		request.Rr = subDomain
		request.Marker = pageMarker
		request.MaxKeys = pageSize

		response, err := d.client.ListRecord(zoneName, request)
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

func (d *DNSProvider) addOrUpdateDNSRecord(zoneName, subDomain, value string) error {
	record, err := d.findDNSRecord(zoneName, subDomain)
	if err != nil {
		return err
	}

	if record == nil {
		request := &bcedns.CreateRecordRequest{
			Type:  "TXT",
			Rr:    subDomain,
			Value: value,
			Ttl:   &d.config.TTL,
		}
		err := d.client.CreateRecord(zoneName, request, d.generateClientToken())
		return err
	} else {
		request := &bcedns.UpdateRecordRequest{
			Type:  "TXT",
			Rr:    subDomain,
			Value: value,
			Ttl:   &d.config.TTL,
		}
		err := d.client.UpdateRecord(zoneName, record.Id, request, d.generateClientToken())
		return err
	}
}

func (d *DNSProvider) removeDNSRecord(zoneName, subDomain string) error {
	record, err := d.findDNSRecord(zoneName, subDomain)
	if err != nil {
		return err
	}

	if record == nil {
		return nil
	} else {
		err = d.client.DeleteRecord(zoneName, record.Id, d.generateClientToken())
		return err
	}
}

func (d *DNSProvider) generateClientToken() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
