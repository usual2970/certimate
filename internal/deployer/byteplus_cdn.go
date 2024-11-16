package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	bytepluscdn "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/byteplus-cdn"

	"github.com/byteplus-sdk/byteplus-sdk-golang/service/cdn"
	xerrors "github.com/pkg/errors"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
)

type ByteplusCDNDeployer struct {
	option      *DeployerOption
	infos       []string
	sdkClient   *cdn.CDN
	sslUploader uploader.Uploader
}

func NewByteplusCDNDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.ByteplusAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}
	client := cdn.NewInstance()
	client.Client.SetAccessKey(access.AccessKey)
	client.Client.SetSecretKey(access.SecretKey)
	uploader, err := bytepluscdn.New(&bytepluscdn.ByteplusCDNUploaderConfig{
		AccessKey: access.AccessKey,
		SecretKey: access.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}
	return &ByteplusCDNDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *ByteplusCDNDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *ByteplusCDNDeployer) GetInfos() []string {
	return d.infos
}

func (d *ByteplusCDNDeployer) Deploy(ctx context.Context) error {
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
		// 获取证书可以部署的域名
		// REF: https://docs.byteplus.com/en/docs/byteplus-cdn/reference-describecertconfig-9ea17
		describeCertConfigReq := &cdn.DescribeCertConfigRequest{
			CertId: upres.CertId,
		}
		describeCertConfigResp, err := d.sdkClient.DescribeCertConfig(describeCertConfigReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'cdn.DescribeCertConfig'")
		}
		for i := range describeCertConfigResp.Result.CertNotConfig {
			// 当前未启用 HTTPS 的加速域名列表。
			domains = append(domains, describeCertConfigResp.Result.CertNotConfig[i].Domain)
		}
		for i := range describeCertConfigResp.Result.OtherCertConfig {
			// 已启用了 HTTPS 的加速域名列表。这些加速域名关联的证书不是您指定的证书。
			domains = append(domains, describeCertConfigResp.Result.OtherCertConfig[i].Domain)
		}
		for i := range describeCertConfigResp.Result.SpecifiedCertConfig {
			// 已启用了 HTTPS 的加速域名列表。这些加速域名关联了您指定的证书。
			d.infos = append(d.infos, fmt.Sprintf("%s域名已配置该证书", describeCertConfigResp.Result.SpecifiedCertConfig[i].Domain))
		}
		if len(domains) == 0 {
			if len(describeCertConfigResp.Result.SpecifiedCertConfig) > 0 {
				// 所有匹配的域名都配置了该证书，跳过部署
				return nil
			} else {
				return xerrors.Errorf("未查询到匹配的域名: %s", configDomain)
			}
		}
	} else {
		domains = append(domains, configDomain)
	}
	// 部署证书
	// REF: https://github.com/byteplus-sdk/byteplus-sdk-golang/blob/master/service/cdn/api_list.go#L306
	for i := range domains {
		batchDeployCertReq := &cdn.BatchDeployCertRequest{
			CertId: upres.CertId,
			Domain: domains[i],
		}
		batchDeployCertResp, err := d.sdkClient.BatchDeployCert(batchDeployCertReq)
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'cdn.BatchDeployCert'")
		} else {
			d.infos = append(d.infos, toStr(fmt.Sprintf("%s域名的证书已修改", domains[i]), batchDeployCertResp))
		}
	}

	return nil
}
