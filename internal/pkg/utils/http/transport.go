package httputil

import (
	"net"
	"net/http"
	"time"
)

// 创建并返回一个 [http.DefaultTransport] 对象副本。
//
// 出参：
//   - transport: [http.DefaultTransport] 对象副本。
func NewDefaultTransport() *http.Transport {
	if http.DefaultTransport != nil {
		if t, ok := http.DefaultTransport.(*http.Transport); ok {
			return t.Clone()
		}
	}

	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
