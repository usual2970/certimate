package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	xerrors "github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	tcTeo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploaderTcSsl "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/tencentcloud-ssl"
)

type TencentTEODeployer struct {
	option *DeployerOption
	infos  []string

	sdkClients  *tencentTEODeployerSdkClients
	sslUploader uploader.Uploader
}

type tencentTEODeployerSdkClients struct {
	ssl *tcSsl.Client
	teo *tcTeo.Client
}

func NewTencentTEODeployer(option *DeployerOption) (Deployer, error) {
	access := &domain.TencentAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, xerrors.Wrap(err, "failed to get access")
	}

	clients, err := (&TencentTEODeployer{}).createSdkClients(
		access.SecretId,
		access.SecretKey,
	)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk clients")
	}

	uploader, err := uploaderTcSsl.New(&uploaderTcSsl.TencentCloudSSLUploaderConfig{
		SecretId:  access.SecretId,
		SecretKey: access.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &TencentTEODeployer{
		option:      option,
		infos:       make([]string, 0),
		sdkClients:  clients,
		sslUploader: uploader,
	}, nil
}

func (d *TencentTEODeployer) GetID() string {
	return fmt.Sprintf("%s-%s", d.option.AccessRecord.GetString("name"), d.option.AccessRecord.Id)
}

func (d *TencentTEODeployer) GetInfos() []string {
	return d.infos
}

func (d *TencentTEODeployer) Deploy(ctx context.Context) error {
	tcZoneId := d.option.DeployConfig.GetConfigAsString("zoneId")
	if tcZoneId == "" {
		return xerrors.New("`zoneId` is required")
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, d.option.Certificate.Certificate, d.option.Certificate.PrivateKey)
	if err != nil {
		return err
	}

	d.infos = append(d.infos, toStr("已上传证书", upres))

	// 配置域名证书
	// REF: https://cloud.tencent.com/document/product/1552/80764
	modifyHostsCertificateReq := tcTeo.NewModifyHostsCertificateRequest()
	modifyHostsCertificateReq.ZoneId = common.StringPtr(tcZoneId)
	modifyHostsCertificateReq.Mode = common.StringPtr("sslcert")
	modifyHostsCertificateReq.Hosts = common.StringPtrs(strings.Split(strings.ReplaceAll(d.option.Domain, "\r\n", "\n"), "\n"))
	modifyHostsCertificateReq.ServerCertInfo = []*tcTeo.ServerCertInfo{{CertId: common.StringPtr(upres.CertId)}}
	modifyHostsCertificateResp, err := d.sdkClients.teo.ModifyHostsCertificate(modifyHostsCertificateReq)
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'teo.ModifyHostsCertificate'")
	}

	d.infos = append(d.infos, toStr("已配置域名证书", modifyHostsCertificateResp.Response))

	return nil
}

func (d *TencentTEODeployer) createSdkClients(secretId, secretKey string) (*tencentTEODeployerSdkClients, error) {
	credential := common.NewCredential(secretId, secretKey)

	sslClient, err := tcSsl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	teoClient, err := tcTeo.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &tencentTEODeployerSdkClients{
		ssl: sslClient,
		teo: teoClient,
	}, nil
}
