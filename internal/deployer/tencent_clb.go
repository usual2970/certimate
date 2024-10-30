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

type TencentCLBDeployer struct {
	option *DeployerOption
	infos  []string

	sdkClients  *tencentCLBDeployerSdkClients
	sslUploader uploader.Uploader
}

type tencentCLBDeployerSdkClients struct {
	ssl *tcSsl.Client
}

func NewTencentCLBDeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.TencentAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	clients, err := (&TencentCLBDeployer{}).createSdkClients(
		access.SecretId,
		access.SecretKey,
		option.DeployConfig.GetConfigAsString("region"),
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

	return &TencentCLBDeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClients:  clients,
		sslUploader: uploader,
	}, nil
}

func (d *TencentCLBDeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *TencentCLBDeployer) GetInfo() []string {
	return d.infos
}

func (d *TencentCLBDeployer) Deploy(ctx context.Context) error {
	// TODO: 直接部署方式

	// 通过 SSL 服务部署到云资源实例
	err := d.deployToInstanceUseSsl(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (d *TencentCLBDeployer) createSdkClients(secretId, secretKey, region string) (*tencentCLBDeployerSdkClients, error) {
	credential := common.NewCredential(secretId, secretKey)

	sslClient, err := tcSsl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &tencentCLBDeployerSdkClients{
		ssl: sslClient,
	}, nil
}

func (d *TencentCLBDeployer) deployToInstanceUseSsl(ctx context.Context) error {
	tcLoadbalancerId := d.option.DeployConfig.GetConfigAsString("clbId")
	tcListenerId := d.option.DeployConfig.GetConfigAsString("lsnId")
	tcDomain := d.option.DeployConfig.GetConfigAsString("domain")
	if tcLoadbalancerId == "" {
		return errors.New("`loadbalancerId` is required")
	}
	if tcListenerId == "" {
		return errors.New("`listenerId` is required")
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 证书部署到 CLB 实例
	// REF: https://cloud.tencent.com/document/product/400/91667
	deployCertificateInstanceReq := tcSsl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(upres.CertId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("clb")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	if tcDomain == "" {
		// 未开启 SNI，只需指定到监听器
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s|%s", tcLoadbalancerId, tcListenerId)})
	} else {
		// 开启 SNI，需指定到域名（支持泛域名）
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s|%s|%s", tcLoadbalancerId, tcListenerId, tcDomain)})
	}
	deployCertificateInstanceResp, err := d.sdkClients.ssl.DeployCertificateInstance(deployCertificateInstanceReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'ssl.DeployCertificateInstance'")
	}

	d.infos = append(d.infos, toStr("已部署证书到云资源实例", deployCertificateInstanceResp.Response))

	return nil
}
