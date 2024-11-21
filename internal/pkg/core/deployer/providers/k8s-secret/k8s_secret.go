package k8ssecret

import (
	"context"
	"errors"
	"strings"

	xerrors "github.com/pkg/errors"
	k8sCore "k8s.io/api/core/v1"
	k8sMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/utils/x509"
)

type K8sSecretDeployerConfig struct {
	// kubeconfig 文件内容。
	KubeConfig string `json:"kubeConfig,omitempty"`
	// K8s 命名空间。
	Namespace string `json:"namespace,omitempty"`
	// K8s Secret 名称。
	SecretName string `json:"secretName"`
	// K8s Secret 中用于存放证书的 Key。
	SecretDataKeyForCrt string `json:"secretDataKeyForCrt,omitempty"`
	// K8s Secret 中用于存放私钥的 Key。
	SecretDataKeyForKey string `json:"secretDataKeyForKey,omitempty"`
}

type K8sSecretDeployer struct {
	config *K8sSecretDeployerConfig
	logger deployer.Logger
}

var _ deployer.Deployer = (*K8sSecretDeployer)(nil)

func New(config *K8sSecretDeployerConfig) (*K8sSecretDeployer, error) {
	return NewWithLogger(config, deployer.NewNilLogger())
}

func NewWithLogger(config *K8sSecretDeployerConfig, logger deployer.Logger) (*K8sSecretDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	return &K8sSecretDeployer{
		logger: logger,
		config: config,
	}, nil
}

func (d *K8sSecretDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	if d.config.Namespace == "" {
		return nil, errors.New("config `namespace` is required")
	}
	if d.config.SecretName == "" {
		return nil, errors.New("config `secretName` is required")
	}
	if d.config.SecretDataKeyForCrt == "" {
		return nil, errors.New("config `secretDataKeyForCrt` is required")
	}
	if d.config.SecretDataKeyForKey == "" {
		return nil, errors.New("config `secretDataKeyForKey` is required")
	}

	certX509, err := x509.ParseCertificateFromPEM(certPem)
	if err != nil {
		return nil, err
	}

	// 连接
	client, err := createK8sClient(d.config.KubeConfig)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create k8s client")
	}

	var secretPayload *k8sCore.Secret
	secretAnnotations := map[string]string{
		"certimate/common-name":       certX509.Subject.CommonName,
		"certimate/subject-sn":        certX509.Subject.SerialNumber,
		"certimate/subject-alt-names": strings.Join(certX509.DNSNames, ","),
		"certimate/issuer-sn":         certX509.Issuer.SerialNumber,
		"certimate/issuer-org":        strings.Join(certX509.Issuer.Organization, ","),
	}

	// 获取 Secret 实例，如果不存在则创建
	secretPayload, err = client.CoreV1().Secrets(d.config.Namespace).Get(context.TODO(), d.config.SecretName, k8sMeta.GetOptions{})
	if err != nil {
		secretPayload = &k8sCore.Secret{
			TypeMeta: k8sMeta.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: k8sMeta.ObjectMeta{
				Name:        d.config.SecretName,
				Annotations: secretAnnotations,
			},
			Type: k8sCore.SecretType("kubernetes.io/tls"),
		}
		secretPayload.Data = make(map[string][]byte)
		secretPayload.Data[d.config.SecretDataKeyForCrt] = []byte(certPem)
		secretPayload.Data[d.config.SecretDataKeyForKey] = []byte(privkeyPem)

		_, err = client.CoreV1().Secrets(d.config.Namespace).Create(context.TODO(), secretPayload, k8sMeta.CreateOptions{})
		if err != nil {
			return nil, xerrors.Wrap(err, "failed to create k8s secret")
		} else {
			d.logger.Logf("k8s secret created", secretPayload)
			return &deployer.DeployResult{}, nil
		}
	}

	// 更新 Secret 实例
	secretPayload.Type = k8sCore.SecretType("kubernetes.io/tls")
	if secretPayload.ObjectMeta.Annotations == nil {
		secretPayload.ObjectMeta.Annotations = secretAnnotations
	} else {
		for k, v := range secretAnnotations {
			secretPayload.ObjectMeta.Annotations[k] = v
		}
	}
	secretPayload, err = client.CoreV1().Secrets(d.config.Namespace).Update(context.TODO(), secretPayload, k8sMeta.UpdateOptions{})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to update k8s secret")
	}

	d.logger.Logf("k8s secret updated", secretPayload)

	return &deployer.DeployResult{}, nil
}

func createK8sClient(kubeConfig string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if kubeConfig == "" {
		config, err = rest.InClusterConfig()
	} else {
		kubeConfig, err := clientcmd.NewClientConfigFromBytes([]byte(kubeConfig))
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
