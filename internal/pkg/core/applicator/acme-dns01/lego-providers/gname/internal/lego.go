package internal

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"

	gnamesdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/gname"
)

const (
	envNamespace = "GNAME_"

	EnvAppID  = envNamespace + "APP_ID"
	EnvAppKey = envNamespace + "APP_KEY"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvHTTPTimeout        = envNamespace + "HTTP_TIMEOUT"
)

var _ challenge.ProviderTimeout = (*DNSProvider)(nil)

type Config struct {
	AppID  string
	AppKey string

	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int
	HTTPTimeout        time.Duration
}

type DNSProvider struct {
	client *gnamesdk.Client
	config *Config
}

func NewDefaultConfig() *Config {
	return &Config{
		TTL:                env.GetOrDefaultInt(EnvTTL, 300),
		PropagationTimeout: env.GetOrDefaultSecond(EnvPropagationTimeout, 2*time.Minute),
		PollingInterval:    env.GetOrDefaultSecond(EnvPollingInterval, dns01.DefaultPollingInterval),
		HTTPTimeout:        env.GetOrDefaultSecond(EnvHTTPTimeout, 30*time.Second),
	}
}

func NewDNSProvider() (*DNSProvider, error) {
	values, err := env.Get(EnvAppID, EnvAppKey)
	if err != nil {
		return nil, fmt.Errorf("gname: %w", err)
	}

	config := NewDefaultConfig()
	config.AppID = values[EnvAppID]
	config.AppKey = values[EnvAppKey]

	return NewDNSProviderConfig(config)
}

func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("gname: the configuration of the DNS provider is nil")
	}

	client := gnamesdk.NewClient(config.AppID, config.AppKey).
		WithTimeout(config.HTTPTimeout)

	return &DNSProvider{
		client: client,
		config: config,
	}, nil
}

func (d *DNSProvider) Present(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("gname: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("gname: %w", err)
	}

	if err := d.addOrUpdateDNSRecord(dns01.UnFqdn(authZone), subDomain, info.Value); err != nil {
		return fmt.Errorf("gname: %w", err)
	}

	return nil
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("gname: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("gname: %w", err)
	}

	if err := d.removeDNSRecord(dns01.UnFqdn(authZone), subDomain); err != nil {
		return fmt.Errorf("gname: %w", err)
	}

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}

func (d *DNSProvider) findDNSRecord(zoneName, subDomain string) (*gnamesdk.ResolutionRecord, error) {
	page := int32(1)
	pageSize := int32(20)
	for {
		request := &gnamesdk.ListDomainResolutionRequest{}
		request.ZoneName = zoneName
		request.Page = &page
		request.PageSize = &pageSize

		response, err := d.client.ListDomainResolution(request)
		if err != nil {
			return nil, err
		}

		for _, record := range response.Data {
			if record.RecordType == "TXT" && record.RecordName == subDomain {
				return record, nil
			}
		}

		if len(response.Data) == 0 {
			break
		}
		if response.Page*response.PageSize >= response.Count {
			break
		}

		page++
	}

	return nil, nil
}

func (d *DNSProvider) addOrUpdateDNSRecord(zoneName, subDomain, value string) error {
	record, err := d.findDNSRecord(zoneName, subDomain)
	if err != nil {
		return err
	}

	if record == nil {
		request := &gnamesdk.AddDomainResolutionRequest{
			ZoneName:    zoneName,
			RecordType:  "TXT",
			RecordName:  subDomain,
			RecordValue: value,
			TTL:         int32(d.config.TTL),
		}
		_, err := d.client.AddDomainResolution(request)
		return err
	} else {
		recordId, _ := record.ID.Int64()
		request := &gnamesdk.ModifyDomainResolutionRequest{
			ID:          recordId,
			ZoneName:    zoneName,
			RecordType:  "TXT",
			RecordName:  subDomain,
			RecordValue: value,
			TTL:         int32(d.config.TTL),
		}
		_, err := d.client.ModifyDomainResolution(request)
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
	}

	recordId, _ := record.ID.Int64()
	request := &gnamesdk.DeleteDomainResolutionRequest{
		ZoneName: zoneName,
		RecordID: recordId,
	}
	_, err = d.client.DeleteDomainResolution(request)
	return err
}
