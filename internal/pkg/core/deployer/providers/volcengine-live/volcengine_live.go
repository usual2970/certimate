package volcenginelive

import (
	"context"
	"errors"
	"fmt"
	"strings"

	xerrors "github.com/pkg/errors"
	veLive "github.com/volcengine/volc-sdk-golang/service/live/v20230101"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploaderp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/volcengine-live"
)

type VolcEngineLiveDeployerConfig struct {
	// 火山引擎 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 火山引擎 AccessKeySecret。
	AccessKeySecret string `json:"accessKeySecret"`
	// 直播流域名（支持泛域名）。
	Domain string `json:"domain"`
}

type VolcEngineLiveDeployer struct {
	config      *VolcEngineLiveDeployerConfig
	logger      logger.Logger
	sdkClient   *veLive.Live
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*VolcEngineLiveDeployer)(nil)

func New(config *VolcEngineLiveDeployerConfig) (*VolcEngineLiveDeployer, error) {
	return NewWithLogger(config, logger.NewNilLogger())
}

func NewWithLogger(config *VolcEngineLiveDeployerConfig, logger logger.Logger) (*VolcEngineLiveDeployer, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	client := veLive.NewInstance()
	client.SetAccessKey(config.AccessKeyId)
	client.SetSecretKey(config.AccessKeySecret)

	uploader, err := uploaderp.New(&uploaderp.VolcEngineLiveUploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &VolcEngineLiveDeployer{
		logger:      logger,
		config:      config,
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *VolcEngineLiveDeployer) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 Live
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	domains := make([]string, 0)
	if strings.HasPrefix(d.config.Domain, "*.") {
		listDomainDetailPageNum := int32(1)
		listDomainDetailPageSize := int32(1000)
		listDomainDetailTotal := 0
		for {
			// 查询域名列表
			// REF: https://www.volcengine.com/docs/6469/1186277#%E6%9F%A5%E8%AF%A2%E5%9F%9F%E5%90%8D%E5%88%97%E8%A1%A8
			listDomainDetailReq := &veLive.ListDomainDetailBody{
				PageNum:  listDomainDetailPageNum,
				PageSize: listDomainDetailPageSize,
			}
			listDomainDetailResp, err := d.sdkClient.ListDomainDetail(ctx, listDomainDetailReq)
			if err != nil {
				return nil, xerrors.Wrap(err, "failed to execute sdk request 'live.ListDomainDetail'")
			}

			if listDomainDetailResp.Result.DomainList != nil {
				for _, item := range listDomainDetailResp.Result.DomainList {
					// 仅匹配泛域名的下一级子域名
					wildcardDomain := strings.TrimPrefix(d.config.Domain, "*")
					if strings.HasSuffix(item.Domain, wildcardDomain) && !strings.Contains(strings.TrimSuffix(item.Domain, wildcardDomain), ".") {
						domains = append(domains, item.Domain)
					}
				}
			}

			listDomainDetailLen := len(listDomainDetailResp.Result.DomainList)
			if listDomainDetailLen < int(listDomainDetailPageSize) || int(listDomainDetailResp.Result.Total) <= listDomainDetailTotal+listDomainDetailLen {
				break
			} else {
				listDomainDetailPageNum++
				listDomainDetailTotal += listDomainDetailLen
			}
		}

		if len(domains) == 0 {
			return nil, xerrors.Errorf("未查询到匹配的域名: %s", d.config.Domain)
		}
	} else {
		domains = append(domains, d.config.Domain)
	}

	if len(domains) > 0 {
		var errs []error

		for _, domain := range domains {
			// 绑定证书
			// REF: https://www.volcengine.com/docs/6469/1186278#%E7%BB%91%E5%AE%9A%E8%AF%81%E4%B9%A6
			bindCertReq := &veLive.BindCertBody{
				ChainID: upres.CertId,
				Domain:  domain,
				HTTPS:   ve.Bool(true),
			}
			bindCertResp, err := d.sdkClient.BindCert(ctx, bindCertReq)
			if err != nil {
				errs = append(errs, err)
			} else {
				d.logger.Logt(fmt.Sprintf("已绑定证书到域名 %s", domain), bindCertResp)
			}
		}

		if len(errs) > 0 {
			return nil, errors.Join(errs...)
		}
	}

	return &deployer.DeployResult{}, nil
}
