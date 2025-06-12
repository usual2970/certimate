package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"
	"gitlab.ecloud.com/ecloud/ecloudsdkclouddns"
	"gitlab.ecloud.com/ecloud/ecloudsdkclouddns/model"
	"gitlab.ecloud.com/ecloud/ecloudsdkcore/config"
)

const (
	envNamespace = "CMCCCLOUD_"

	EnvAccessKey = envNamespace + "ACCESS_KEY"
	EnvSecretKey = envNamespace + "SECRET_KEY"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvReadTimeOut        = envNamespace + "READ_TIMEOUT"
	EnvConnectTimeout     = envNamespace + "CONNECT_TIMEOUT"
)

var _ challenge.ProviderTimeout = (*DNSProvider)(nil)

type Config struct {
	AccessKey string
	SecretKey string

	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int32
	ReadTimeOut        int
	ConnectTimeout     int
}

type DNSProvider struct {
	client *ecloudsdkclouddns.Client
	config *Config
}

func NewDefaultConfig() *Config {
	return &Config{
		ReadTimeOut:        env.GetOrDefaultInt(EnvReadTimeOut, 30),
		ConnectTimeout:     env.GetOrDefaultInt(EnvConnectTimeout, 30),
		TTL:                int32(env.GetOrDefaultInt(EnvTTL, 600)),
		PropagationTimeout: env.GetOrDefaultSecond(EnvPropagationTimeout, 2*time.Minute),
		PollingInterval:    env.GetOrDefaultSecond(EnvPollingInterval, dns01.DefaultPollingInterval),
	}
}

func NewDNSProvider() (*DNSProvider, error) {
	values, err := env.Get(EnvAccessKey, EnvSecretKey)
	if err != nil {
		return nil, fmt.Errorf("cmccecloud: %w", err)
	}

	cfg := NewDefaultConfig()
	cfg.AccessKey = values[EnvAccessKey]
	cfg.SecretKey = values[EnvSecretKey]

	return NewDNSProviderConfig(cfg)
}

func NewDNSProviderConfig(cfg *Config) (*DNSProvider, error) {
	if cfg == nil {
		return nil, errors.New("cmccecloud: the configuration of the DNS provider is nil")
	}

	client := ecloudsdkclouddns.NewClient(&config.Config{
		AccessKey: cfg.AccessKey,
		SecretKey: cfg.SecretKey,
		// 资源池常量见: https://ecloud.10086.cn/op-help-center/doc/article/54462
		// 默认全局
		PoolId:         "CIDC-CORE-00",
		ReadTimeOut:    cfg.ReadTimeOut,
		ConnectTimeout: cfg.ConnectTimeout,
	})

	return &DNSProvider{
		client: client,
		config: cfg,
	}, nil
}

func (d *DNSProvider) Present(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	zoneName, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("cmccecloud: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, zoneName)
	if err != nil {
		return fmt.Errorf("cmccecloud: %w", err)
	}

	readDomain := strings.Trim(zoneName, ".")
	record, err := d.getDomainRecord(readDomain, subDomain)
	if err != nil {
		return err
	}

	if record == nil {
		resp, err := d.client.CreateRecordOpenapi(&model.CreateRecordOpenapiRequest{
			CreateRecordOpenapiBody: &model.CreateRecordOpenapiBody{
				LineId:      "0", // 默认线路
				Rr:          subDomain,
				DomainName:  readDomain,
				Description: "certimate acme",
				Type:        model.CreateRecordOpenapiBodyTypeEnumTxt,
				Value:       info.Value,
				Ttl:         &d.config.TTL,
			},
		})
		if err != nil {
			return fmt.Errorf("cmccecloud: %w", err)
		}

		if resp.State != model.CreateRecordOpenapiResponseStateEnumOk {
			return fmt.Errorf("cmccecloud: create record failed, response state: %s, message: %s, code: %s", resp.State, resp.ErrorMessage, resp.ErrorCode)
		}

		return nil
	} else {
		resp, err := d.client.ModifyRecordOpenapi(&model.ModifyRecordOpenapiRequest{
			ModifyRecordOpenapiBody: &model.ModifyRecordOpenapiBody{
				RecordId:    record.RecordId,
				Rr:          subDomain,
				DomainName:  readDomain,
				Description: "certmate acme",
				LineId:      "0",
				Type:        model.ModifyRecordOpenapiBodyTypeEnumTxt,
				Value:       info.Value,
				Ttl:         &d.config.TTL,
			},
		})
		if err != nil {
			return fmt.Errorf("cmccecloud: %w", err)
		}

		if resp.State != model.ModifyRecordOpenapiResponseStateEnumOk {
			return fmt.Errorf("cmccecloud: create record failed, response state: %s", resp.State)
		}

		return nil
	}
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	challengeInfo := dns01.GetChallengeInfo(domain, keyAuth)

	zoneName, err := dns01.FindZoneByFqdn(challengeInfo.FQDN)
	if err != nil {
		return fmt.Errorf("cmccecloud: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(challengeInfo.FQDN, zoneName)
	if err != nil {
		return fmt.Errorf("cmccecloud: %w", err)
	}

	readDomain := strings.Trim(zoneName, ".")
	record, err := d.getDomainRecord(readDomain, subDomain)
	if err != nil {
		return err
	}

	if record == nil {
		return nil
	} else {
		resp, err := d.client.DeleteRecordOpenapi(&model.DeleteRecordOpenapiRequest{
			DeleteRecordOpenapiBody: &model.DeleteRecordOpenapiBody{
				RecordIdList: []string{record.RecordId},
			},
		})
		if err != nil {
			return fmt.Errorf("cmccecloud: %w", err)
		}
		if resp.State != model.DeleteRecordOpenapiResponseStateEnumOk {
			return fmt.Errorf("cmccecloud: delete record failed, unexpected response state: %s", resp.State)
		}
	}

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}

func (d *DNSProvider) getDomainRecord(domain string, rr string) (*model.ListRecordOpenapiResponseData, error) {
	pageSize := int32(50)
	page := int32(1)
	for {
		resp, err := d.client.ListRecordOpenapi(&model.ListRecordOpenapiRequest{
			ListRecordOpenapiBody: &model.ListRecordOpenapiBody{
				DomainName: domain,
			},
			ListRecordOpenapiQuery: &model.ListRecordOpenapiQuery{
				PageSize: &pageSize,
				Page:     &page,
			},
		})
		if err != nil {
			return nil, err
		}
		if resp.State != model.ListRecordOpenapiResponseStateEnumOk {
			respStr, _ := json.Marshal(resp)
			return nil, fmt.Errorf("cmccecloud: request error: %s", string(respStr))
		}

		if resp.Body.Data != nil {
			for _, item := range *resp.Body.Data {
				if item.Rr == rr {
					return &item, nil
				}
			}
		}

		if resp.Body.TotalPages == nil || page >= *resp.Body.TotalPages {
			return nil, nil
		}

		page++
	}
}
