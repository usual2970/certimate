package deployer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	xerrors "github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
)

type TencentCOSDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClient   *tcSsl.Client
	sslUploader uploader.Uploader
}

func NewTencentCOSDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.TencentAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	client, err := (&TencentCOSDeployer{}).createSdkClient(
		access.SecretId,
		access.SecretKey,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk clients")
	}

	uploader, err := uploader.NewTencentCloudSSLUploader(&uploader.TencentCloudSSLUploaderConfig{
		SecretId:  access.SecretId,
		SecretKey: access.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &TencentCOSDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *TencentCOSDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *TencentCOSDeployer) GetInfos() []string {
	return d.infos
}

func (d *TencentCOSDeployer) Deploy(ctx context.Context) error {
	tcRegion := d.option.DeployConfig.GetConfigAsString("region")
	tcBucket := d.option.DeployConfig.GetConfigAsString("bucket")
	tcDomain := d.option.DeployConfig.GetConfigAsString("domain")
	if tcBucket == "" {
		return errors.New("`bucket` is required")
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 证书部署到 COS 实例
	// REF: https://cloud.tencent.com/document/product/400/91667
	deployCertificateInstanceReq := tcSsl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(upres.CertId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("cos")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s#%s#%s", tcRegion, tcBucket, tcDomain)})
	deployCertificateInstanceResp, err := d.sdkClient.DeployCertificateInstance(deployCertificateInstanceReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'ssl.DeployCertificateInstance'")
	}

	d.infos = append(d.infos, toStr("已部署证书到云资源实例", deployCertificateInstanceResp.Response))

	return nil
}

func (d *TencentCOSDeployer) createSdkClient(secretId, secretKey string) (*tcSsl.Client, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tcSsl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return client, nil
}
