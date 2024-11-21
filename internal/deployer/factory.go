package deployer

import (
	"encoding/json"
	"fmt"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	providerAliyunOss "github.com/usual2970/certimate/internal/pkg/core/deployer/providers/aliyun-oss"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

// TODO: 该方法目前未实际使用，将在后续迭代中替换
func createDeployer(target string, accessConfig string, deployConfig map[string]any) (deployer.Deployer, deployer.Logger, error) {
	logger := deployer.NewDefaultLogger()

	switch target {
	case targetAliyunOSS:
		access := &domain.AliyunAccess{}
		if err := json.Unmarshal([]byte(accessConfig), access); err != nil {
			return nil, nil, err
		}

		deployer, err := providerAliyunOss.NewWithLogger(&providerAliyunOss.AliyunOSSDeployerConfig{
			AccessKeyId:     access.AccessKeyId,
			AccessKeySecret: access.AccessKeySecret,
			Region:          maps.GetValueAsString(deployConfig, "region"),
			Bucket:          maps.GetValueAsString(deployConfig, "bucket"),
			Domain:          maps.GetValueAsString(deployConfig, "domain"),
		}, logger)
		return deployer, logger, err
	}

	return nil, nil, fmt.Errorf("unsupported deployer target: %s", target)
}
