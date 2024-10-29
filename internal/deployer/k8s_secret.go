package deployer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	xerrors "github.com/pkg/errors"
	k8sCore "k8s.io/api/core/v1"
	k8sMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
)

type K8sSecretDeployer struct {
	option *DeployerOption
	infos  []string

	k8sClient *kubernetes.Clientset
}

func NewK8sSecretDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.KubernetesAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	client, err := (&K8sSecretDeployer{}).createK8sClient(access)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create k8s client")
	}

	return &K8sSecretDeployer{
		option:    option,
		infos:     make([]string, 0),
		k8sClient: client,
	}, nil
}

func (d *K8sSecretDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *K8sSecretDeployer) GetInfo() []string {
	return d.infos
}

func (d *K8sSecretDeployer) Deploy(ctx context.Context) error {
	namespace := d.option.DeployConfig.GetConfigAsString("namespace")
	secretName := d.option.DeployConfig.GetConfigAsString("secretName")
	secretDataKeyForCrt := d.option.DeployConfig.GetConfigOrDefaultAsString("secretDataKeyForCrt", "tls.crt")
	secretDataKeyForKey := d.option.DeployConfig.GetConfigOrDefaultAsString("secretDataKeyForKey", "tls.key")
	if namespace == "" {
		namespace = "default"
	}
	if secretName == "" {
		return errors.New("`secretName` is required")
	}

	certX509, err := x509.ParseCertificateFromPEM(d.option.Certificate.Certificate)
	if err != nil {
		return err
	}

	secretPayload := k8sCore.Secret{
		TypeMeta: k8sMeta.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: k8sMeta.ObjectMeta{
			Name: secretName,
			Annotations: map[string]string{
				"certimate/domains":             d.option.Domain,
				"certimate/alt-names":           strings.Join(certX509.DNSNames, ","),
				"certimate/common-name":         certX509.Subject.CommonName,
				"certimate/issuer-organization": strings.Join(certX509.Issuer.Organization, ","),
			},
		},
		Type: k8sCore.SecretType("kubernetes.io/tls"),
	}
	secretPayload.Data = make(map[string][]byte)
	secretPayload.Data[secretDataKeyForCrt] = []byte(d.option.Certificate.Certificate)
	secretPayload.Data[secretDataKeyForKey] = []byte(d.option.Certificate.PrivateKey)

	// 获取 Secret 实例
	_, err = d.k8sClient.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, k8sMeta.GetOptions{})
	if err != nil {
		_, err = d.k8sClient.CoreV1().Secrets(namespace).Create(context.TODO(), &secretPayload, k8sMeta.CreateOptions{})
		if err != nil {
			return xerrors.Wrap(err, "failed to create k8s secret")
		} else {
			d.infos = append(d.infos, toStr("Certificate has been created in K8s Secret", nil))
			return nil
		}
	}

	// 更新 Secret 实例
	_, err = d.k8sClient.CoreV1().Secrets(namespace).Update(context.TODO(), &secretPayload, k8sMeta.UpdateOptions{})
	if err != nil {
		return xerrors.Wrap(err, "failed to update k8s secret")
	}

	d.infos = append(d.infos, toStr("Certificate has been updated to K8s Secret", nil))

	return nil
}

func (d *K8sSecretDeployer) createK8sClient(access *domain.KubernetesAccess) (*kubernetes.Clientset, error) {
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
