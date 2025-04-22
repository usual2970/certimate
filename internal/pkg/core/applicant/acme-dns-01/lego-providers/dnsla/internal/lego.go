package lego_dnsla

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"

	dnslasdk "github.com/usual2970/certimate/internal/pkg/sdk3rd/dnsla"
)

const (
	envNamespace = "DNSLA_"

	EnvAPIId     = envNamespace + "API_ID"
	EnvAPISecret = envNamespace + "API_KEY"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvHTTPTimeout        = envNamespace + "HTTP_TIMEOUT"
)

var _ challenge.ProviderTimeout = (*DNSProvider)(nil)

type Config struct {
	APIId     string
	APISecret string

	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int
	HTTPTimeout        time.Duration
}

type DNSProvider struct {
	client *dnslasdk.Client
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
	values, err := env.Get(EnvAPIId, EnvAPISecret)
	if err != nil {
		return nil, fmt.Errorf("dnsla: %w", err)
	}

	config := NewDefaultConfig()
	config.APIId = values[EnvAPIId]
	config.APISecret = values[EnvAPISecret]

	return NewDNSProviderConfig(config)
}

func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("dnsla: the configuration of the DNS provider is nil")
	}

	client := dnslasdk.NewClient(config.APIId, config.APISecret).
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
		return fmt.Errorf("dnsla: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("dnsla: %w", err)
	}

	if err := d.addOrUpdateDNSRecord(dns01.UnFqdn(authZone), subDomain, info.Value); err != nil {
		return fmt.Errorf("dnsla: %w", err)
	}

	return nil
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("dnsla: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("dnsla: %w", err)
	}

	if err := d.removeDNSRecord(dns01.UnFqdn(authZone), subDomain); err != nil {
		return fmt.Errorf("dnsla: %w", err)
	}

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}

func (d *DNSProvider) getDNSZone(zoneName string) (*dnslasdk.DomainInfo, error) {
	pageIndex := 1
	pageSize := 100
	for {
		request := &dnslasdk.ListDomainsRequest{
			PageIndex: int32(pageIndex),
			PageSize:  int32(pageSize),
		}
		response, err := d.client.ListDomains(request)
		if err != nil {
			return nil, err
		}

		if response.Data != nil {
			for _, item := range response.Data.Results {
				if strings.TrimRight(item.Domain, ".") == zoneName || strings.TrimRight(item.DisplayDomain, ".") == zoneName {
					return item, nil
				}
			}
		}

		if response.Data == nil || len(response.Data.Results) < pageSize {
			break
		}

		pageIndex++
	}

	return nil, fmt.Errorf("dnsla: zone %s not found", zoneName)
}

func (d *DNSProvider) getDNSZoneAndRecord(zoneName, subDomain string) (*dnslasdk.DomainInfo, *dnslasdk.RecordInfo, error) {
	zone, err := d.getDNSZone(zoneName)
	if err != nil {
		return nil, nil, err
	}

	pageIndex := 1
	pageSize := 100
	for {
		request := &dnslasdk.ListRecordsRequest{
			DomainId:  zone.Id,
			Host:      &subDomain,
			PageIndex: int32(pageIndex),
			PageSize:  int32(pageSize),
		}
		response, err := d.client.ListRecords(request)
		if err != nil {
			return zone, nil, err
		}

		if response.Data != nil {
			for _, record := range response.Data.Results {
				if record.Type == 16 && (record.Host == subDomain || record.DisplayHost == subDomain) {
					return zone, record, nil
				}
			}
		}

		if response.Data == nil || len(response.Data.Results) < pageSize {
			break
		}

		pageIndex++
	}

	return zone, nil, nil
}

func (d *DNSProvider) addOrUpdateDNSRecord(zoneName, subDomain, value string) error {
	zone, record, err := d.getDNSZoneAndRecord(zoneName, subDomain)
	if err != nil {
		return err
	}

	if record == nil {
		request := &dnslasdk.CreateRecordRequest{
			DomainId: zone.Id,
			Type:     16,
			Host:     subDomain,
			Data:     value,
			Ttl:      int32(d.config.TTL),
		}
		_, err := d.client.CreateRecord(request)
		return err
	} else {
		reqType := int32(16)
		reqTtl := int32(d.config.TTL)
		request := &dnslasdk.UpdateRecordRequest{
			Id:   record.Id,
			Type: &reqType,
			Host: &subDomain,
			Data: &value,
			Ttl:  &reqTtl,
		}
		_, err := d.client.UpdateRecord(request)
		return err
	}
}

func (d *DNSProvider) removeDNSRecord(zoneName, subDomain string) error {
	_, record, err := d.getDNSZoneAndRecord(zoneName, subDomain)
	if err != nil {
		return err
	}

	if record == nil {
		return nil
	} else {
		request := &dnslasdk.DeleteRecordRequest{
			Id: record.Id,
		}
		_, err = d.client.DeleteRecord(request)
		return err
	}
}
