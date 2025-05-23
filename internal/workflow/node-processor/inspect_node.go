package nodeprocessor

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/usual2970/certimate/internal/domain"
)

type inspectNode struct {
	node *domain.WorkflowNode
	*nodeProcessor
	*nodeOutputer
}

func NewInspectNode(node *domain.WorkflowNode) *inspectNode {
	return &inspectNode{
		node:          node,
		nodeProcessor: newNodeProcessor(node),
		nodeOutputer:  newNodeOutputer(),
	}
}

func (n *inspectNode) Process(ctx context.Context) error {
	n.logger.Info("entering inspect certificate node...")

	nodeConfig := n.node.GetConfigForInspect()

	err := n.inspect(ctx, nodeConfig)
	if err != nil {
		n.logger.Warn("inspect certificate failed: " + err.Error())
		return err
	}

	return nil
}

func (n *inspectNode) inspect(ctx context.Context, nodeConfig domain.WorkflowNodeConfigForInspect) error {
	maxRetries := 3
	retryInterval := 2 * time.Second

	var lastError error
	var certInfo *x509.Certificate

	host := nodeConfig.Host

	port := nodeConfig.Port
	if port == "" {
		port = "443"
	}

	domain := nodeConfig.Domain
	if domain == "" {
		domain = host
	}

	path := nodeConfig.Path
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	targetAddr := fmt.Sprintf("%s:%s", host, port)
	n.logger.Info(fmt.Sprintf("Inspecting certificate at %s (validating domain: %s)", targetAddr, domain))

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			n.logger.Info(fmt.Sprintf("Retry #%d connecting to %s", attempt, targetAddr))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryInterval):
				// Wait for retry interval
			}
		}

		transport := &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         domain, // Set SNI to domain for proper certificate selection
			},
			ForceAttemptHTTP2: false,
			DisableKeepAlives: true,
		}

		client := &http.Client{
			Transport: transport,
			Timeout:   15 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		scheme := "https"
		urlStr := fmt.Sprintf("%s://%s", scheme, targetAddr)
		if path != "" {
			urlStr = urlStr + path
		}

		req, err := http.NewRequestWithContext(ctx, "HEAD", urlStr, nil)
		if err != nil {
			lastError = fmt.Errorf("failed to create HTTP request: %w", err)
			n.logger.Warn(fmt.Sprintf("Request creation attempt #%d failed: %s", attempt+1, lastError.Error()))
			continue
		}

		if domain != host {
			req.Host = domain
		}

		req.Header.Set("User-Agent", "CertificateValidator/1.0")
		req.Header.Set("Accept", "*/*")

		resp, err := client.Do(req)
		if err != nil {
			lastError = fmt.Errorf("HTTP request failed: %w", err)
			n.logger.Warn(fmt.Sprintf("Connection attempt #%d failed: %s", attempt+1, lastError.Error()))
			continue
		}

		if resp.TLS == nil || len(resp.TLS.PeerCertificates) == 0 {
			resp.Body.Close()
			lastError = fmt.Errorf("no TLS certificates received in HTTP response")
			n.logger.Warn(fmt.Sprintf("Certificate retrieval attempt #%d failed: %s", attempt+1, lastError.Error()))
			continue
		}

		certInfo = resp.TLS.PeerCertificates[0]
		resp.Body.Close()

		lastError = nil
		n.logger.Info(fmt.Sprintf("Successfully retrieved certificate from %s", targetAddr))
		break
	}

	if lastError != nil {
		return fmt.Errorf("failed to retrieve certificate after %d attempts: %w", maxRetries, lastError)
	}

	if certInfo == nil {
		outputs := map[string]any{
			outputCertificateValidatedKey: "false",
			outputCertificateDaysLeftKey:  "0",
		}
		n.setOutputs(outputs)
		return nil
	}

	now := time.Now()

	isValidTime := now.Before(certInfo.NotAfter) && now.After(certInfo.NotBefore)

	domainMatch := true
	if err := certInfo.VerifyHostname(domain); err != nil {
		domainMatch = false
	}

	isValid := isValidTime && domainMatch

	daysRemaining := math.Floor(certInfo.NotAfter.Sub(now).Hours() / 24)

	isValidStr := "false"
	if isValid {
		isValidStr = "true"
	}

	outputs := map[string]any{
		outputCertificateValidatedKey: isValidStr,
		outputCertificateDaysLeftKey:  fmt.Sprintf("%d", int(daysRemaining)),
	}

	n.setOutputs(outputs)

	n.logger.Info(fmt.Sprintf("Certificate inspection completed - Target: %s, Domain: %s, Valid: %s, Days Remaining: %d",
		targetAddr, domain, isValidStr, int(daysRemaining)))

	return nil
}

func (n *inspectNode) setOutputs(outputs map[string]any) {
	n.outputs = outputs
}
