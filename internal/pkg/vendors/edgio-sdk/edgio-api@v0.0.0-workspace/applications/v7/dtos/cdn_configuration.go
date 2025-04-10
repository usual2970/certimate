package dtos

import "encoding/json"

type CDNConfiguration struct {
	ConfigurationID        string             `json:"id"`
	EnvironmentID          string             `json:"environment_id"`
	Rules                  json.RawMessage    `json:"rules"`
	Origins                []Origin           `json:"origins"`
	Hostnames              []Hostname         `json:"hostnames"`
	Experiments            *[]string          `json:"experiments,omitempty"`
	EdgeFunctionsSources   *map[string]string `json:"edge_functions_sources,omitempty"`
	EdgeFunctionInitScript *string            `json:"edge_function_init_script,omitempty"`
}

type Origin struct {
	Name                string     `json:"name"`
	Type                *string    `json:"type,omitempty"`
	Hosts               []Host     `json:"hosts"`
	Balancer            *string    `json:"balancer,omitempty"`
	OverrideHostHeader  *string    `json:"override_host_header,omitempty"`
	Shields             *Shields   `json:"shields,omitempty"`
	PciCertifiedShields *bool      `json:"pci_certified_shields,omitempty"`
	TLSVerify           *TLSVerify `json:"tls_verify,omitempty"`
	Retry               *Retry     `json:"retry,omitempty"`
}

type Host struct {
	Weight                   *int64      `json:"weight,omitempty"`
	DNSMaxTTL                *int64      `json:"dns_max_ttl,omitempty"`
	DNSPreference            *string     `json:"dns_preference,omitempty"`
	MaxHardPool              *int64      `json:"max_hard_pool,omitempty"`
	DNSMinTTL                *int64      `json:"dns_min_ttl,omitempty"`
	Location                 *[]Location `json:"location,omitempty"`
	MaxPool                  *int64      `json:"max_pool,omitempty"`
	Balancer                 *string     `json:"balancer,omitempty"`
	Scheme                   *string     `json:"scheme,omitempty"`
	OverrideHostHeader       *string     `json:"override_host_header,omitempty"`
	SNIHintAndStrictSanCheck *string     `json:"sni_hint_and_strict_san_check,omitempty"`
	UseSNI                   *bool       `json:"use_sni,omitempty"`
}

type Location struct {
	Port     *int64  `json:"port,omitempty"`
	Hostname *string `json:"hostname,omitempty"`
}

type Shields struct {
	Apac   *string `json:"apac,omitempty"`
	Emea   *string `json:"emea,omitempty"`
	USWest *string `json:"us_west,omitempty"`
	USEast *string `json:"us_east,omitempty"`
}

type TLSVerify struct {
	UseSNI                   *bool     `json:"use_sni,omitempty"`
	SNIHintAndStrictSanCheck *string   `json:"sni_hint_and_strict_san_check,omitempty"`
	AllowSelfSignedCerts     *bool     `json:"allow_self_signed_certs,omitempty"`
	PinnedCerts              *[]string `json:"pinned_certs,omitempty"`
}

type Retry struct {
	StatusCodes            *[]int64 `json:"status_codes,omitempty"`
	IgnoreRetryAfterHeader *bool    `json:"ignore_retry_after_header,omitempty"`
	AfterSeconds           *int64   `json:"after_seconds,omitempty"`
	MaxRequests            *int64   `json:"max_requests,omitempty"`
	MaxWaitSeconds         *int64   `json:"max_wait_seconds,omitempty"`
}

type Hostname struct {
	Hostname          *string `json:"hostname,omitempty"`
	DefaultOriginName *string `json:"default_origin_name,omitempty"`
	ReportCode        *int64  `json:"report_code,omitempty"`
	TLS               *TLS    `json:"tls,omitempty"`
	Directory         *string `json:"directory,omitempty"`
}

type TLS struct {
	NPN                 *bool   `json:"npn,omitempty"`
	ALPN                *bool   `json:"alpn,omitempty"`
	Protocols           *string `json:"protocols,omitempty"`
	UseSigAlgs          *bool   `json:"use_sigalgs,omitempty"`
	SNI                 *bool   `json:"sni,omitempty"`
	SniStrict           *bool   `json:"sni_strict,omitempty"`
	SniHostMatch        *bool   `json:"sni_host_match,omitempty"`
	ClientRenegotiation *bool   `json:"client_renegotiation,omitempty"`
	Options             *string `json:"options,omitempty"`
	CipherList          *string `json:"cipher_list,omitempty"`
	NamedCurve          *string `json:"named_curve,omitempty"`
	OCSP                *bool   `json:"oscp,omitempty"`
	PEM                 *string `json:"pem,omitempty"`
	CA                  *string `json:"ca,omitempty"`
}
