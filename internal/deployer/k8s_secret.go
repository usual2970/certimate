package deployer

import (
	"context"
	"encoding/json"
	"fmt"

	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"certimate/internal/domain"
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
	return fmt.Sprintf("%s-%s", d.option.AceessRecord.GetString("name"), d.option.AceessRecord.Id)
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

	d.infos = append(d.infos, toStr("kubeClient 创建成功", nil))

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

	// 获取 Secret 实例
	secret, err := client.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, k8sMetaV1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get k8s secret: %w", err)
	}

	// 更新 Secret Data
	secret.Data[secretDataKeyForCrt] = []byte(d.option.Certificate.Certificate)
	secret.Data[secretDataKeyForKey] = []byte(d.option.Certificate.PrivateKey)
	_, err = client.CoreV1().Secrets(namespace).Update(context.TODO(), secret, k8sMetaV1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update k8s secret: %w", err)
	}

	d.infos = append(d.infos, toStr("证书已更新到 K8s Secret", nil))

	return nil
}

func (d *K8sSecretDeployer) createClient(access *domain.KubernetesAccess) (*kubernetes.Clientset, error) {
	kubeConfig, err := clientcmd.Load([]byte(access.KubeConfig))
	if err != nil {
		return nil, err
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: ""},
		&clientcmd.ConfigOverrides{CurrentContext: kubeConfig.CurrentContext},
	)
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
