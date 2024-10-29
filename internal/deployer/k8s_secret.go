package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
)

type K8sSecretDeployer struct {
	option *DeployerOption
	infos  []string
}

func NewK8sSecretDeployer(option *DeployerOption) (Deployer, error) {
	return &K8sSecretDeployer{
		option: option,
		infos:  make([]string, 0),
	}, nil
}

func (d *K8sSecretDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *K8sSecretDeployer) GetInfo() []string {
	return d.infos
}

func (d *K8sSecretDeployer) Deploy(ctx context.Context) error {
	access := &domain.KubernetesAccess{}
	if err := json.Unmarshal([]byte(d.option.Access), access); err != nil {
		return err
	}

	client, err := d.createClient(access)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("kubeClient create success.", nil))

	namespace := getDeployString(d.option.DeployConfig, "namespace")
	if namespace == "" {
		namespace = "default"
	}

	secretName := getDeployString(d.option.DeployConfig, "secretName")
	if secretName == "" {
		return fmt.Errorf("k8s secret name is empty")
	}

	secretDataKeyForCrt := getDeployString(d.option.DeployConfig, "secretDataKeyForCrt")
	if secretDataKeyForCrt == "" {
		namespace = "tls.crt"
	}

	secretDataKeyForKey := getDeployString(d.option.DeployConfig, "secretDataKeyForKey")
	if secretDataKeyForKey == "" {
		namespace = "tls.key"
	}

	certificate, err := x509.ParseCertificateFromPEM(d.option.Certificate.Certificate)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	secretPayload := corev1.Secret{
		TypeMeta: k8sMetaV1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name: secretName,
			Annotations: map[string]string{
				"certimate/domains":             d.option.Domain,
				"certimate/alt-names":           strings.Join(certificate.DNSNames, ","),
				"certimate/common-name":         certificate.Subject.CommonName,
				"certimate/issuer-organization": strings.Join(certificate.Issuer.Organization, ","),
			},
		},
		Type: corev1.SecretType("kubernetes.io/tls"),
	}

	secretPayload.Data = make(map[string][]byte)
	secretPayload.Data[secretDataKeyForCrt] = []byte(d.option.Certificate.Certificate)
	secretPayload.Data[secretDataKeyForKey] = []byte(d.option.Certificate.PrivateKey)

	// 获取 Secret 实例
	_, err = client.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, k8sMetaV1.GetOptions{})
	if err != nil {
		_, err = client.CoreV1().Secrets(namespace).Create(context.TODO(), &secretPayload, k8sMetaV1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create k8s secret: %w", err)
		} else {
			d.infos = append(d.infos, toStr("Certificate has been created in K8s Secret", nil))
			return nil
		}
	}

	// 更新 Secret 实例
	_, err = client.CoreV1().Secrets(namespace).Update(context.TODO(), &secretPayload, k8sMetaV1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update k8s secret: %w", err)
	}

	d.infos = append(d.infos, toStr("Certificate has been updated to K8s Secret", nil))

	return nil
}

func (d *K8sSecretDeployer) createClient(access *domain.KubernetesAccess) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if access.KubeConfig == "" {
		config, err = rest.InClusterConfig()
	} else {
		kubeConfig, err := clientcmd.NewClientConfigFromBytes([]byte(access.KubeConfig))
		if err != nil {
			return nil, err
		}
		config, err = kubeConfig.ClientConfig()

	}
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
