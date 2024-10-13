package deployer

import (
	"certimate/internal/applicant"
	"certimate/internal/domain"
	"certimate/internal/utils/rand"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	cas20200407 "github.com/alibabacloud-go/cas-20200407/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type aliyun struct {
	client *cas20200407.Client
	option *DeployerOption
	infos  []string
}

func NewAliyun(option *DeployerOption) (Deployer, error) {
	access := &domain.AliyunAccess{}
	json.Unmarshal([]byte(option.Access), access)
	a := &aliyun{
		option: option,
		infos:  make([]string, 0),
	}
	client, err := a.createClient(access.AccessKeyId, access.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	a.client = client
	return a, nil

}

func (a *aliyun) GetID() string {
	return fmt.Sprintf("%s-%s", a.option.AceessRecord.GetString("name"), a.option.AceessRecord.Id)
}

func (a *aliyun) GetInfo() []string {
	return a.infos
}

func (a *aliyun) Deploy(ctx context.Context) error {

	// 查询有没有对应的资源
	resource, err := a.resource()
	if err != nil {
		return err
	}

	a.infos = append(a.infos, toStr("查询对应的资源", resource))

	// 查询有没有对应的联系人
	contacts, err := a.contacts()
	if err != nil {
		return err
	}

	a.infos = append(a.infos, toStr("查询联系人", contacts))

	// 上传证书
	certId, err := a.uploadCert(&a.option.Certificate)
	if err != nil {
		return err
	}

	a.infos = append(a.infos, toStr("上传证书", certId))

	// 部署证书
	jobId, err := a.deploy(resource, certId, contacts)
	if err != nil {
		return err
	}

	a.infos = append(a.infos, toStr("创建部署证书任务", jobId))

	// 等待部署成功
	err = a.updateDeployStatus(*jobId)
	if err != nil {
		return err
	}

	// 部署成功后删除旧的证书
	a.deleteCert(resource)

	return nil
}

func (a *aliyun) updateDeployStatus(jobId int64) error {
	// 查询部署状态
	req := &cas20200407.UpdateDeploymentJobStatusRequest{
		JobId: tea.Int64(jobId),
	}

	resp, err := a.client.UpdateDeploymentJobStatus(req)
	if err != nil {
		return err
	}
	a.infos = append(a.infos, toStr("查询对应的资源", resp))
	return nil
}

func (a *aliyun) deleteCert(resource *cas20200407.ListCloudResourcesResponseBodyData) error {
	// 查询有没有对应的资源
	if resource.CertId == nil {
		return nil
	}

	// 删除证书
	_, err := a.client.DeleteUserCertificate(&cas20200407.DeleteUserCertificateRequest{
		CertId: resource.CertId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *aliyun) contacts() ([]*cas20200407.ListContactResponseBodyContactList, error) {
	listContactRequest := &cas20200407.ListContactRequest{}
	runtime := &util.RuntimeOptions{}

	resp, err := a.client.ListContactWithOptions(listContactRequest, runtime)
	if err != nil {
		return nil, err
	}
	if resp.Body.TotalCount == nil {
		return nil, errors.New("no contact found")
	}

	return resp.Body.ContactList, nil
}

func (a *aliyun) deploy(resource *cas20200407.ListCloudResourcesResponseBodyData, certId int64, contacts []*cas20200407.ListContactResponseBodyContactList) (*int64, error) {
	contactIds := make([]string, 0, len(contacts))
	for _, contact := range contacts {
		contactIds = append(contactIds, fmt.Sprintf("%d", *contact.ContactId))
	}
	// 部署证书
	createCloudResourceRequest := &cas20200407.CreateDeploymentJobRequest{
		CertIds:     tea.String(fmt.Sprintf("%d", certId)),
		Name:        tea.String(a.option.Domain + rand.RandStr(6)),
		JobType:     tea.String("user"),
		ResourceIds: tea.String(fmt.Sprintf("%d", *resource.Id)),
		ContactIds:  tea.String(strings.Join(contactIds, ",")),
	}
	runtime := &util.RuntimeOptions{}

	resp, err := a.client.CreateDeploymentJobWithOptions(createCloudResourceRequest, runtime)
	if err != nil {
		return nil, err
	}
	return resp.Body.JobId, nil
}

func (a *aliyun) uploadCert(cert *applicant.Certificate) (int64, error) {
	uploadUserCertificateRequest := &cas20200407.UploadUserCertificateRequest{
		Cert: &cert.Certificate,
		Key:  &cert.PrivateKey,
		Name: tea.String(a.option.Domain + rand.RandStr(6)),
	}
	runtime := &util.RuntimeOptions{}

	resp, err := a.client.UploadUserCertificateWithOptions(uploadUserCertificateRequest, runtime)
	if err != nil {
		return 0, err
	}

	return *resp.Body.CertId, nil
}

func (a *aliyun) createClient(accessKeyId, accessKeySecret string) (_result *cas20200407.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	config.Endpoint = tea.String("cas.aliyuncs.com")
	_result = &cas20200407.Client{}
	_result, _err = cas20200407.NewClient(config)
	return _result, _err
}

func (a *aliyun) resource() (*cas20200407.ListCloudResourcesResponseBodyData, error) {

	listCloudResourcesRequest := &cas20200407.ListCloudResourcesRequest{
		CloudProduct: tea.String(a.option.Product),
		Keyword:      tea.String(getDeployString(a.option.DeployConfig, "domain")),
	}

	resp, err := a.client.ListCloudResources(listCloudResourcesRequest)
	if err != nil {
		return nil, err
	}

	if *resp.Body.Total == 0 {
		return nil, errors.New("no resource found")
	}

	return resp.Body.Data[0], nil
}
