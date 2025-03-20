package lego_dynv6

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"
	"github.com/libdns/dynv6"
	"github.com/libdns/libdns"
)

const (
	envNamespace = "DYNV6_"

	EnvHTTPToken = envNamespace + "HTTP_TOKEN"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
)

var _ challenge.ProviderTimeout = (*DNSProvider)(nil)

type Config struct {
	HTTPToken string

	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int
}

type DNSProvider struct {
	client *dynv6.Provider
	config *Config
}

func NewDefaultConfig() *Config {
	return &Config{
		TTL:                env.GetOrDefaultInt(EnvTTL, dns01.DefaultTTL),
		PropagationTimeout: env.GetOrDefaultSecond(EnvPropagationTimeout, 2*time.Minute),
		PollingInterval:    env.GetOrDefaultSecond(EnvPollingInterval, dns01.DefaultPollingInterval),
	}
}

func NewDNSProvider() (*DNSProvider, error) {
	values, err := env.Get(EnvHTTPToken)
	if err != nil {
		return nil, fmt.Errorf("dynv6: %w", err)
	}

	config := NewDefaultConfig()
	config.HTTPToken = values[EnvHTTPToken]

	return NewDNSProviderConfig(config)
}

func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("dynv6: the configuration of the DNS provider is nil")
	}

	client := &dynv6.Provider{Token: config.HTTPToken}

	return &DNSProvider{
		client: client,
		config: config,
	}, nil
}

func (d *DNSProvider) Present(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("dynv6: %w", err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("dynv6: %w", err)
	}

	if err := d.addOrUpdateDNSRecord(dns01.UnFqdn(authZone), subDomain, info.Value); err != nil {
		return fmt.Errorf("dynv6: %w", err)
	}

	return nil
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("dynv6: %w", err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("dynv6: %w", err)
	}

	if err := d.removeDNSRecord(dns01.UnFqdn(authZone), subDomain); err != nil {
		return fmt.Errorf("dynv6: %w", err)
	}

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}

func (d *DNSProvider) getDNSRecord(zoneName, subDomain string) (*libdns.Record, error) {
	records, err := d.client.GetRecords(context.Background(), zoneName)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if record.Type == "TXT" && record.Name == subDomain {
			return &record, nil
		}
	}

	return nil, nil
}

func (d *DNSProvider) addOrUpdateDNSRecord(zoneName, subDomain, value string) error {
	record, err := d.getDNSRecord(zoneName, subDomain)
	if err != nil {
		return err
	}

	if record == nil {
		record = &libdns.Record{
			Type:  "TXT",
			Name:  subDomain,
			Value: value,
			TTL:   time.Duration(d.config.TTL) * time.Second,
		}
		_, err := d.client.AppendRecords(context.Background(), zoneName, []libdns.Record{*record})
		return err
	} else {
		record.Value = value
		_, err := d.client.SetRecords(context.Background(), zoneName, []libdns.Record{*record})
		return err
	}
}

func (d *DNSProvider) removeDNSRecord(zoneName, subDomain string) error {
	record, err := d.getDNSRecord(zoneName, subDomain)
	if err != nil {
		return err
	}

	if record == nil {
		return nil
	} else {
		_, err = d.client.DeleteRecords(context.Background(), zoneName, []libdns.Record{*record})
		return err
	}
}
