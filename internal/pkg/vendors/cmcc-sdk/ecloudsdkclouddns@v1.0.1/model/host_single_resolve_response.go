// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package model



type HostSingleResolveResponse struct {

	// 该域名的IPv6解析结果，是一个列表，可能包括0个、1个或多个IP地址；仅当query=6时返回这个字段。
	Ipsv6 []string `json:"ipsv6,omitempty"`

	// 该域名原始TTL，即权威NS上配置的域名TTL值。
	OriginTtl *int32 `json:"origin_ttl,omitempty"`

	// 解析的域名
	Host string `json:"host,omitempty"`

	// 该域名的IPv4解析结果，是一个列表，可能包括0个、1个或多个IP地址；仅当query=4时返回这个字段。
	Ips []string `json:"ips,omitempty"`

	// 该域名解析结果的TTL缓存时间。
	Ttl *int32 `json:"ttl,omitempty"`
}
