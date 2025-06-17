package internal

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"

	ctyundns "github.com/certimate-go/certimate/pkg/sdk3rd/ctyun/dns"
	xtypes "github.com/certimate-go/certimate/pkg/utils/types"
)

const (
	envNamespace = "CTYUNSMARTDNS_"

	EnvAccessKeyID     = envNamespace + "ACCESS_KEY_ID"
	EnvSecretAccessKey = envNamespace + "SECRET_ACCESS_KEY"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvHTTPTimeout        = envNamespace + "HTTP_TIMEOUT"
)

var _ challenge.ProviderTimeout = (*DNSProvider)(nil)

type Config struct {
	AccessKeyId     string
	SecretAccessKey string

	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int
	HTTPTimeout        time.Duration
}

type DNSProvider struct {
	client *ctyundns.Client
	config *Config
}

func NewDefaultConfig() *Config {
	return &Config{
		TTL:                env.GetOrDefaultInt(EnvTTL, 600),
		PropagationTimeout: env.GetOrDefaultSecond(EnvPropagationTimeout, 2*time.Minute),
		HTTPTimeout:        env.GetOrDefaultSecond(EnvHTTPTimeout, 30*time.Second),
	}
}

func NewDNSProvider() (*DNSProvider, error) {
	values, err := env.Get(EnvAccessKeyID, EnvSecretAccessKey)
	if err != nil {
		return nil, fmt.Errorf("ctyun: %w", err)
	}

	config := NewDefaultConfig()
	config.AccessKeyId = values[EnvAccessKeyID]
	config.SecretAccessKey = values[EnvSecretAccessKey]

	return NewDNSProviderConfig(config)
}

func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("ctyun: the configuration of the DNS provider is nil")
	}

	client, err := ctyundns.NewClient(config.AccessKeyId, config.SecretAccessKey)
	if err != nil {
		return nil, err
	} else {
		client.SetTimeout(config.HTTPTimeout)
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
		return fmt.Errorf("ctyun: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("ctyun: %w", err)
	}

	if err := d.addOrUpdateDNSRecord(dns01.UnFqdn(authZone), subDomain, info.Value); err != nil {
		return fmt.Errorf("ctyun: %w", err)
	}

	return nil
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("ctyun: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("ctyun: %w", err)
	}

	if err := d.removeDNSRecord(dns01.UnFqdn(authZone), subDomain); err != nil {
		return fmt.Errorf("ctyun: %w", err)
	}

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}

func (d *DNSProvider) findDNSRecordId(zoneName, subDomain string) (int32, error) {
	// 查询解析记录列表
	// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=122&api=11264&data=181&isNormal=1&vid=259
	request := &ctyundns.QueryRecordListRequest{}
	request.Domain = xtypes.ToPtr(zoneName)
	request.Host = xtypes.ToPtr(subDomain)
	request.Type = xtypes.ToPtr("TXT")

	response, err := d.client.QueryRecordList(request)
	if err != nil {
		return 0, err
	}

	if response.ReturnObj == nil || response.ReturnObj.Records == nil || len(response.ReturnObj.Records) == 0 {
		return 0, nil
	}

	return response.ReturnObj.Records[0].RecordId, nil
}

func (d *DNSProvider) addOrUpdateDNSRecord(zoneName, subDomain, value string) error {
	recordId, err := d.findDNSRecordId(zoneName, subDomain)
	if err != nil {
		return err
	}

	if recordId == 0 {
		// 新增解析记录
		// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=122&api=11259&data=181&isNormal=1&vid=259
		request := &ctyundns.AddRecordRequest{
			Domain:   xtypes.ToPtr(zoneName),
			Host:     xtypes.ToPtr(subDomain),
			Type:     xtypes.ToPtr("TXT"),
			LineCode: xtypes.ToPtr("Default"),
			Value:    xtypes.ToPtr(value),
			State:    xtypes.ToPtr(int32(1)),
			TTL:      xtypes.ToPtr(int32(d.config.TTL)),
		}
		_, err := d.client.AddRecord(request)
		return err
	} else {
		// 修改解析记录
		// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=122&api=11261&data=181&isNormal=1&vid=259
		request := &ctyundns.UpdateRecordRequest{
			RecordId: xtypes.ToPtr(recordId),
			Domain:   xtypes.ToPtr(zoneName),
			Host:     xtypes.ToPtr(subDomain),
			Type:     xtypes.ToPtr("TXT"),
			LineCode: xtypes.ToPtr("Default"),
			Value:    xtypes.ToPtr(value),
			State:    xtypes.ToPtr(int32(1)),
			TTL:      xtypes.ToPtr(int32(d.config.TTL)),
		}
		_, err := d.client.UpdateRecord(request)
		return err
	}
}

func (d *DNSProvider) removeDNSRecord(zoneName, subDomain string) error {
	recordId, err := d.findDNSRecordId(zoneName, subDomain)
	if err != nil {
		return err
	}

	if recordId == 0 {
		return nil
	} else {
		// 删除解析记录
		// REF: https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=122&api=11262&data=181&isNormal=1&vid=259
		request := &ctyundns.DeleteRecordRequest{
			RecordId: xtypes.ToPtr(recordId),
		}
		_, err = d.client.DeleteRecord(request)
		return err
	}
}
