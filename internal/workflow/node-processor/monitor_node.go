package nodeprocessor

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/usual2970/certimate/internal/domain"
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
	n.logger.Info("ready to monitor certificate ...")

	nodeConfig := n.node.GetConfigForMonitor()

	targetAddr := fmt.Sprintf("%s:%d", nodeConfig.Host, nodeConfig.Port)
	if nodeConfig.Port == 0 {
		targetAddr = fmt.Sprintf("%s:443", nodeConfig.Host)
	}

	targetDomain := nodeConfig.Domain
	if targetDomain == "" {
		targetDomain = nodeConfig.Host
	}

	n.logger.Info(fmt.Sprintf("retrieving certificate at %s (domain: %s)", targetAddr, targetDomain))

	const MAX_ATTEMPTS = 3
	const RETRY_INTERVAL = 2 * time.Second
	var cert *x509.Certificate
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

		cert, err = n.tryRetrieveCert(ctx, targetAddr, targetDomain, nodeConfig.RequestPath)
		if err == nil {
			break
		}
	}

	if err != nil {
		n.logger.Warn("failed to monitor certificate")
		return err
	} else {
		if cert == nil {
			n.logger.Warn("no ssl certificates retrieved in http response")

			outputs := map[string]any{
				outputCertificateValidatedKey: strconv.FormatBool(false),
				outputCertificateDaysLeftKey:  strconv.FormatInt(0, 10),
			}
			n.setOutputs(outputs)
		} else {
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
			outputs := map[string]any{
				outputCertificateValidatedKey: strconv.FormatBool(validated),
				outputCertificateDaysLeftKey:  strconv.FormatInt(int64(daysLeft), 10),
			}
			n.setOutputs(outputs)

			if validated {
				n.logger.Info(fmt.Sprintf("the certificate is valid, and will expire in %d day(s)", daysLeft))
			} else {
				n.logger.Warn(fmt.Sprintf("the certificate is invalid", validated))
			}
		}
	}

	n.logger.Info("monitoring completed")
	return nil
}

func (n *monitorNode) tryRetrieveCert(ctx context.Context, addr, domain, requestPath string) (_cert *x509.Certificate, _err error) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		ForceAttemptHTTP2: false,
		DisableKeepAlives: true,
		Proxy:             http.ProxyFromEnvironment,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	url := fmt.Sprintf("https://%s/%s", addr, strings.TrimLeft(requestPath, "/"))
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		_err = fmt.Errorf("failed to create http request: %w", err)
		n.logger.Warn(fmt.Sprintf("failed to create http request: %w", err))
		return nil, _err
	}

	req.Header.Set("User-Agent", "certimate")
	resp, err := client.Do(req)
	if err != nil {
		_err = fmt.Errorf("failed to send http request: %w", err)
		n.logger.Warn(fmt.Sprintf("failed to send http request: %w", err))
		return nil, _err
	}
	defer resp.Body.Close()

	if resp.TLS == nil || len(resp.TLS.PeerCertificates) == 0 {
		return nil, _err
	}

	_cert = resp.TLS.PeerCertificates[0]
	return _cert, nil
}

func (n *monitorNode) setOutputs(outputs map[string]any) {
	n.outputs = outputs
}
