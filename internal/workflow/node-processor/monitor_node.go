package nodeprocessor

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/certimate-go/certimate/internal/domain"
	xhttp "github.com/certimate-go/certimate/pkg/utils/http"
)

type monitorNode struct {
	node *domain.WorkflowNode
	*nodeProcessor
	*nodeOutputer
}

func NewMonitorNode(node *domain.WorkflowNode) *monitorNode {
	return &monitorNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),
		nodeOutputer:  newNodeOutputer(),
	}
}

func (n *monitorNode) Process(ctx context.Context) error {
	nodeCfg := n.node.GetConfigForMonitor()
	n.logger.Info("ready to monitor certificate ...", slog.Any("config", nodeCfg))

	targetAddr := net.JoinHostPort(nodeCfg.Host, strconv.Itoa(int(nodeCfg.Port)))
	if nodeCfg.Port == 0 {
		targetAddr = net.JoinHostPort(nodeCfg.Host, "443")
	}

	targetDomain := nodeCfg.Domain
	if targetDomain == "" {
		targetDomain = nodeCfg.Host
	}

	n.logger.Info(fmt.Sprintf("retrieving certificate at %s (domain: %s)", targetAddr, targetDomain))

	const MAX_ATTEMPTS = 3
	const RETRY_INTERVAL = 2 * time.Second
	var certs []*x509.Certificate
	var err error
	for attempt := 0; attempt < MAX_ATTEMPTS; attempt++ {
		if attempt > 0 {
			n.logger.Info(fmt.Sprintf("retry %d time(s) ...", attempt, targetAddr))

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(RETRY_INTERVAL):
			}
		}

		certs, err = n.tryRetrievePeerCertificates(ctx, targetAddr, targetDomain, nodeCfg.RequestPath)
		if err == nil {
			break
		}
	}

	if err != nil {
		n.logger.Warn("failed to monitor certificate")
		return err
	} else {
		if len(certs) == 0 {
			n.logger.Warn("no ssl certificates retrieved in http response")

			n.outputs[outputKeyForCertificateValidity] = strconv.FormatBool(false)
			n.outputs[outputKeyForCertificateDaysLeft] = strconv.FormatInt(0, 10)
		} else {
			cert := certs[0] // 只取证书链中的第一个证书，即服务器证书
			n.logger.Info(fmt.Sprintf("ssl certificate retrieved (serial='%s', subject='%s', issuer='%s', not_before='%s', not_after='%s', sans='%s')",
				cert.SerialNumber, cert.Subject.String(), cert.Issuer.String(),
				cert.NotBefore.Format(time.RFC3339), cert.NotAfter.Format(time.RFC3339),
				strings.Join(cert.DNSNames, ";")),
			)

			now := time.Now()
			isCertPeriodValid := now.Before(cert.NotAfter) && now.After(cert.NotBefore)
			isCertHostMatched := true
			if err := cert.VerifyHostname(targetDomain); err != nil {
				isCertHostMatched = false
			}

			validated := isCertPeriodValid && isCertHostMatched
			daysLeft := int(math.Floor(cert.NotAfter.Sub(now).Hours() / 24))
			n.outputs[outputKeyForCertificateValidity] = strconv.FormatBool(validated)
			n.outputs[outputKeyForCertificateDaysLeft] = strconv.FormatInt(int64(daysLeft), 10)

			if validated {
				n.logger.Info(fmt.Sprintf("the certificate is valid, and will expire in %d day(s)", daysLeft))
			} else {
				if !isCertHostMatched {
					n.logger.Warn("the certificate is invalid, because it is not matched the host")
				} else if !isCertPeriodValid {
					n.logger.Warn("the certificate is invalid, because it is either expired or not yet valid")
				} else {
					n.logger.Warn("the certificate is invalid")
				}
			}
		}
	}

	n.logger.Info("monitoring completed")
	return nil
}

func (n *monitorNode) tryRetrievePeerCertificates(ctx context.Context, addr, domain, requestPath string) ([]*x509.Certificate, error) {
	transport := xhttp.NewDefaultTransport()
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}
	transport.TLSClientConfig.InsecureSkipVerify = true

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	url := fmt.Sprintf("https://%s/%s", addr, strings.TrimLeft(requestPath, "/"))
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		err = fmt.Errorf("failed to create http request: %w", err)
		n.logger.Warn(err.Error())
		return nil, err
	}

	req.Header.Set("User-Agent", "certimate")
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to send http request: %w", err)
		n.logger.Warn(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.TLS == nil || len(resp.TLS.PeerCertificates) == 0 {
		return make([]*x509.Certificate, 0), nil
	}
	return resp.TLS.PeerCertificates, nil
}
