package deployer

import (
	"certimate/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
)

type tencentCdn struct {
	option     *DeployerOption
	credential *common.Credential
	infos      []string
}

func NewTencentCdn(option *DeployerOption) (Deployer, error) {

	access := &domain.TencentAccess{}
	if err := json.Unmarshal([]byte(option.Access), access); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tencent access: %w", err)
	}

	credential := common.NewCredential(
		access.SecretId,
		access.SecretKey,
	)

	return &tencentCdn{
		option:     option,
		credential: credential,
		infos:      make([]string, 0),
	}, nil
}

func (a *tencentCdn) GetID() string {
	return fmt.Sprintf("%s-%s", a.option.AceessRecord.GetString("name"), a.option.AceessRecord.Id)
}

func (t *tencentCdn) GetInfo() []string {
	return t.infos
}

func (t *tencentCdn) Deploy(ctx context.Context) error {

	// 查询有没有对应的资源
	resource, err := t.resource()
	if err != nil {
		return fmt.Errorf("failed to get resource: %w", err)
	}

	t.infos = append(t.infos, toStr("查询对应的资源", resource))

	// 上传证书
	certId, err := t.uploadCert()
	if err != nil {
		return fmt.Errorf("failed to upload certificate: %w", err)
	}
	t.infos = append(t.infos, toStr("上传证书", certId))

	if err := t.deploy(resource, certId); err != nil {
		return fmt.Errorf("failed to deploy: %w", err)
	}

	return nil
}

func (t *tencentCdn) uploadCert() (string, error) {

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ssl.tencentcloudapi.com"

	client, _ := ssl.NewClient(t.credential, "", cpf)

	request := ssl.NewUploadCertificateRequest()

	request.CertificatePublicKey = common.StringPtr(t.option.Certificate.Certificate)
	request.CertificatePrivateKey = common.StringPtr(t.option.Certificate.PrivateKey)
	request.Alias = common.StringPtr(t.option.Domain)
	request.Repeatable = common.BoolPtr(true)

	response, err := client.UploadCertificate(request)
	if err != nil {
		return "", fmt.Errorf("failed to upload certificate: %w", err)
	}

	return *response.Response.CertificateId, nil
}

func (t *tencentCdn) deploy(resource *tag.ResourceTagMapping, certId string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ssl.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := ssl.NewClient(t.credential, "", cpf)

	resourceId, err := getResourceId(resource)
	if err != nil {
		return fmt.Errorf("failed to get resource id: %w", err)
	}

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := ssl.NewDeployCertificateInstanceRequest()

	request.CertificateId = common.StringPtr(certId)
	request.InstanceIdList = common.StringPtrs([]string{resourceId})
	request.ResourceType = common.StringPtr("cdn")
	request.Status = common.Int64Ptr(1)

	// 返回的resp是一个DeployCertificateInstanceResponse的实例，与请求对象对应
	resp, err := client.DeployCertificateInstance(request)

	if err != nil {
		return fmt.Errorf("failed to deploy certificate: %w", err)
	}
	t.infos = append(t.infos, toStr("部署证书", resp.Response))
	return nil
}

func (t *tencentCdn) resource() (*tag.ResourceTagMapping, error) {
	request := tag.NewGetResourcesRequest()
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tag.tencentcloudapi.com"

	client, err := tag.NewClient(t.credential, "", cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	response, err := client.GetResources(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get resources: %w", err)
	}

	for _, resource := range response.Response.ResourceTagMappingList {
		if t.compare(resource) {
			return resource, nil
		}
	}

	return nil, errors.New("no resource found")

}

func (t *tencentCdn) compare(resource *tag.ResourceTagMapping) bool {
	slices := strings.Split(*resource.Resource, "/")
	if len(slices) != 3 {
		return false
	}

	typeSlices := strings.Split(slices[0], "::")
	if len(typeSlices) != 3 {
		return false
	}

	if typeSlices[1] != "cdn" || slices[2] != t.option.Domain {
		return false
	}

	return true

}

func getResourceId(resource *tag.ResourceTagMapping) (string, error) {
	slices := strings.Split(*resource.Resource, "/")
	if len(slices) != 3 {
		return "", errors.New("invalid resource")
	}
	return slices[2], nil
}
