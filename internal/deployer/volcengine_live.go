package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	xerrors "github.com/pkg/errors"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	volcenginelive "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/volcengine-live"
	"github.com/usual2970/certimate/internal/pkg/utils/cast"
	"github.com/volcengine/volc-sdk-golang/base"
	live "github.com/volcengine/volc-sdk-golang/service/live/v20230101"
)

type VolcengineLiveDeployer struct {
	option      *DeployerOption
	infos       []string
	sdkClient   *live.Live
	sslUploader uploader.Uploader
}

func NewVolcengineLiveDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.VolcengineAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}
	client := live.NewInstance()
	client.SetCredential(base.Credentials{
		AccessKeyID:     access.AccessKeyID,
		SecretAccessKey: access.SecretAccessKey,
	})
	uploader, err := volcenginelive.New(&volcenginelive.VolcengineLiveUploaderConfig{
		AccessKeyId:     access.AccessKeyID,
		AccessKeySecret: access.SecretAccessKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}
	return &VolcengineLiveDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *VolcengineLiveDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *VolcengineLiveDeployer) GetInfos() []string {
	return d.infos
}

func (d *VolcengineLiveDeployer) Deploy(ctx context.Context) error {
	apiCtx := context.Background()
	// 上传证书
	upres, err := d.sslUploader.Upload(apiCtx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	domains := make([]string, 0)
	configDomain := d.option.DeployConfig.GetConfigAsString("domain")
	if strings.HasPrefix(configDomain, "*.") {
		// 如果是泛域名，获取所有的域名并匹配
		matchDomains, err := d.getDomainsByWildcardDomain(apiCtx, configDomain)
		if err != nil {
			d.infos = append(d.infos, toStr("获取域名列表失败", upres))
			return xerrors.Wrap(err, "failed to execute sdk request 'live.ListDomainDetail'")
		}
		if len(matchDomains) == 0 {
			return xerrors.Errorf("未查询到匹配的域名: %s", configDomain)
		}
		domains = matchDomains
	} else {
		domains = append(domains, configDomain)
	}

	// 部署证书
	// REF: https://www.volcengine.com/docs/6469/1186278#%E7%BB%91%E5%AE%9A%E8%AF%81%E4%B9%A6d
	for i := range domains {
		bindCertReq := &live.BindCertBody{
			ChainID: upres.CertId,
			Domain:  domains[i],
			HTTPS:   cast.BoolPtr(true),
		}
		bindCertResp, err := d.sdkClient.BindCert(apiCtx, bindCertReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'live.BindCert'")
		} else {
			d.infos = append(d.infos, toStr(fmt.Sprintf("%s域名的证书已修改", domains[i]), bindCertResp))
		}
	}

	return nil
}

func (d *VolcengineLiveDeployer) getDomainsByWildcardDomain(ctx context.Context, wildcardDomain string) ([]string, error) {
	pageNum := int32(1)
	searchTotal := 0
	domains := make([]string, 0)
	for {
		listDomainDetailReq := &live.ListDomainDetailBody{
			PageNum:  pageNum,
			PageSize: 1000,
		}
		// 查询域名列表
		// REF: https://www.volcengine.com/docs/6469/1186277#%E6%9F%A5%E8%AF%A2%E5%9F%9F%E5%90%8D%E5%88%97%E8%A1%A8
		listDomainDetailResp, err := d.sdkClient.ListDomainDetail(ctx, listDomainDetailReq)
		if err != nil {
			return domains, err
		}
		if listDomainDetailResp.Result.DomainList != nil {
			for _, item := range listDomainDetailResp.Result.DomainList {
				if matchWildcardDomain(item.Domain, wildcardDomain) {
					domains = append(domains, item.Domain)
				}
			}
		}
		searchTotal += len(listDomainDetailResp.Result.DomainList)
		if int(listDomainDetailResp.Result.Total) > searchTotal {
			pageNum++
		} else {
			break
		}
	}

	return domains, nil
}

func matchWildcardDomain(domain, wildcardDomain string) bool {
	if strings.HasPrefix(wildcardDomain, "*.") {
		if "*."+domain == wildcardDomain {
			return true
		}
		regexPattern := "^([a-zA-Z0-9_-]+)\\." + regexp.QuoteMeta(wildcardDomain[2:]) + "$"
		regex := regexp.MustCompile(regexPattern)
		return regex.MatchString(domain)
	}
	return domain == wildcardDomain
}
