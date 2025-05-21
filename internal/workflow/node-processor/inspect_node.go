package nodeprocessor

import (
	"context"
	"crypto/tls"
	"fmt"
	"math"
	"net"
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
	n.logger.Info("enter inspect website certificate node ...")

	nodeConfig := n.node.GetConfigForInspect()

	err := n.inspect(ctx, nodeConfig)
	if err != nil {
		n.logger.Warn("inspect website certificate failed: " + err.Error())
		return err
	}

	return nil
}

func (n *inspectNode) inspect(ctx context.Context, nodeConfig domain.WorkflowNodeConfigForInspect) error {
	// 定义重试参数
	maxRetries := 3
	retryInterval := 2 * time.Second

	var cert *tls.Certificate
	var lastError error

	domainWithPort := nodeConfig.Domain + ":" + nodeConfig.Port

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			n.logger.Info(fmt.Sprintf("Retry #%d connecting to %s", attempt, domainWithPort))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryInterval):
				// Wait for retry interval
			}
		}

		dialer := &net.Dialer{
			Timeout: 10 * time.Second,
		}

		conn, err := tls.DialWithDialer(dialer, "tcp", domainWithPort, &tls.Config{
			InsecureSkipVerify: true, // Allow self-signed certificates
		})
		if err != nil {
			lastError = fmt.Errorf("failed to connect to %s: %w", domainWithPort, err)
			n.logger.Warn(fmt.Sprintf("Connection attempt #%d failed: %s", attempt+1, lastError.Error()))
			continue
		}

		// Get certificate information
		certInfo := conn.ConnectionState().PeerCertificates[0]
		conn.Close()

		// Certificate information retrieved successfully
		cert = &tls.Certificate{
			Certificate: [][]byte{certInfo.Raw},
			Leaf:        certInfo,
		}
		lastError = nil
		n.logger.Info(fmt.Sprintf("Successfully retrieved certificate information for %s", domainWithPort))
		break
	}

	if lastError != nil {
		return fmt.Errorf("failed to retrieve certificate after %d attempts: %w", maxRetries, lastError)
	}

	certInfo := cert.Leaf
	now := time.Now()

	isValid := now.Before(certInfo.NotAfter) && now.After(certInfo.NotBefore)

	// Check domain matching
	domainMatch := false
	if len(certInfo.DNSNames) > 0 {
		for _, dnsName := range certInfo.DNSNames {
			if matchDomain(nodeConfig.Domain, dnsName) {
				domainMatch = true
				break
			}
		}
	} else if matchDomain(nodeConfig.Domain, certInfo.Subject.CommonName) {
		domainMatch = true
	}

	isValid = isValid && domainMatch

	daysRemaining := math.Floor(certInfo.NotAfter.Sub(now).Hours() / 24)

	// Set node outputs
	outputs := map[string]any{
		"certificate.validated": isValid,
		"certificate.daysLeft":  daysRemaining,
	}
	n.setOutputs(outputs)

	return nil
}

func (n *inspectNode) setOutputs(outputs map[string]any) {
	n.outputs = outputs
}

func matchDomain(requestDomain, certDomain string) bool {
	if requestDomain == certDomain {
		return true
	}

	if len(certDomain) > 2 && certDomain[0] == '*' && certDomain[1] == '.' {

		wildcardSuffix := certDomain[1:]
		requestDomainLen := len(requestDomain)
		suffixLen := len(wildcardSuffix)

		if requestDomainLen > suffixLen && requestDomain[requestDomainLen-suffixLen:] == wildcardSuffix {
			remainingPart := requestDomain[:requestDomainLen-suffixLen]
			if len(remainingPart) > 0 && !contains(remainingPart, '.') {
				return true
			}
		}
	}

	return false
}

func contains(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}
