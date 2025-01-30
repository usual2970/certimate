// @Title  Golang SDK Client
// @Description  This code is auto generated
// @Author  Ecloud SDK

package ecloudsdkclouddns

import (
    "gitlab.ecloud.com/ecloud/ecloudsdkclouddns/model"
    "gitlab.ecloud.com/ecloud/ecloudsdkcore"
    "gitlab.ecloud.com/ecloud/ecloudsdkcore/config"
)

type Client struct{
    APIClient   *ecloudsdkcore.APIClient
    config      *config.Config
    httpRequest *ecloudsdkcore.HttpRequest
}

func NewClient(config *config.Config) *Client {
    client := &Client{}
    client.config = config
    apiClient := ecloudsdkcore.NewAPIClient()
    httpRequest := ecloudsdkcore.NewDefaultHttpRequest()
    httpRequest.Product = product
    httpRequest.Version = version
    httpRequest.SdkVersion = sdkVersion
    client.httpRequest = httpRequest
    client.APIClient = apiClient
    return client
}

func NewClientByCustomized(config *config.Config,httpRequest *ecloudsdkcore.HttpRequest) *Client {
    client := &Client{}
    client.config = config
    apiClient := ecloudsdkcore.NewAPIClient()
    httpRequest.Product = product
    httpRequest.Version = version
    httpRequest.SdkVersion = sdkVersion
    client.httpRequest = httpRequest
    client.APIClient = apiClient
    return client
}

const (
    product    string = "clouddns"
    version    string = "v1"
    sdkVersion string = "1.0.1"
)

// UnlockDomain 域名解锁
func (c *Client) UnlockDomain(request *model.UnlockDomainRequest) (*model.UnlockDomainResponse, error) {
    c.httpRequest.Action = "unlockDomain"
    c.httpRequest.Body = request
    var returnValue = &model.UnlockDomainResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// LockDomain 域名锁定
func (c *Client) LockDomain(request *model.LockDomainRequest) (*model.LockDomainResponse, error) {
    c.httpRequest.Action = "lockDomain"
    c.httpRequest.Body = request
    var returnValue = &model.LockDomainResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// HostSingleResolve 单域名解析
func (c *Client) HostSingleResolve(request *model.HostSingleResolveRequest) (*model.HostSingleResolveResponse, error) {
    c.httpRequest.Action = "hostSingleResolve"
    c.httpRequest.Body = request
    var returnValue = &model.HostSingleResolveResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ControlRecord 启停解析记录
func (c *Client) ControlRecord(request *model.ControlRecordRequest) (*model.ControlRecordResponse, error) {
    c.httpRequest.Action = "controlRecord"
    c.httpRequest.Body = request
    var returnValue = &model.ControlRecordResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// HostBatchResolve 批量域名解析
func (c *Client) HostBatchResolve(request *model.HostBatchResolveRequest) (*model.HostBatchResolveResponse, error) {
    c.httpRequest.Action = "hostBatchResolve"
    c.httpRequest.Body = request
    var returnValue = &model.HostBatchResolveResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// DomainStatistics 根据域名查询统计量
func (c *Client) DomainStatistics(request *model.DomainStatisticsRequest) (*model.DomainStatisticsResponse, error) {
    c.httpRequest.Action = "domainStatistics"
    c.httpRequest.Body = request
    var returnValue = &model.DomainStatisticsResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// DeleteDomain 删除域名
func (c *Client) DeleteDomain(request *model.DeleteDomainRequest) (*model.DeleteDomainResponse, error) {
    c.httpRequest.Action = "deleteDomain"
    c.httpRequest.Body = request
    var returnValue = &model.DeleteDomainResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// DeleteRecord 删除解析记录
func (c *Client) DeleteRecord(request *model.DeleteRecordRequest) (*model.DeleteRecordResponse, error) {
    c.httpRequest.Action = "deleteRecord"
    c.httpRequest.Body = request
    var returnValue = &model.DeleteRecordResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// CreateDomain 新增域名
func (c *Client) CreateDomain(request *model.CreateDomainRequest) (*model.CreateDomainResponse, error) {
    c.httpRequest.Action = "createDomain"
    c.httpRequest.Body = request
    var returnValue = &model.CreateDomainResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// SingleHostResolve 单域名解析
func (c *Client) SingleHostResolve(request *model.SingleHostResolveRequest) (*model.SingleHostResolveResponse, error) {
    c.httpRequest.Action = "singleHostResolve"
    c.httpRequest.Body = request
    var returnValue = &model.SingleHostResolveResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ControlRecordOpenapi 启停解析记录Openapi
func (c *Client) ControlRecordOpenapi(request *model.ControlRecordOpenapiRequest) (*model.ControlRecordOpenapiResponse, error) {
    c.httpRequest.Action = "controlRecordOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ControlRecordOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// BatchHostResolve 批量域名解析
func (c *Client) BatchHostResolve(request *model.BatchHostResolveRequest) (*model.BatchHostResolveResponse, error) {
    c.httpRequest.Action = "batchHostResolve"
    c.httpRequest.Body = request
    var returnValue = &model.BatchHostResolveResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ModifyRecordOpenapi 修改解析记录Openapi
func (c *Client) ModifyRecordOpenapi(request *model.ModifyRecordOpenapiRequest) (*model.ModifyRecordOpenapiResponse, error) {
    c.httpRequest.Action = "modifyRecordOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ModifyRecordOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// DeleteRecordOpenapi 删除解析记录Openapi
func (c *Client) DeleteRecordOpenapi(request *model.DeleteRecordOpenapiRequest) (*model.DeleteRecordOpenapiResponse, error) {
    c.httpRequest.Action = "deleteRecordOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.DeleteRecordOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListNewRecordLineOpenapi 查询解析线路Openapi
func (c *Client) ListNewRecordLineOpenapi(request *model.ListNewRecordLineOpenapiRequest) (*model.ListNewRecordLineOpenapiResponse, error) {
    c.httpRequest.Action = "listNewRecordLineOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ListNewRecordLineOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// SingleHostTagResolve 单域名解析
func (c *Client) SingleHostTagResolve(request *model.SingleHostTagResolveRequest) (*model.SingleHostTagResolveResponse, error) {
    c.httpRequest.Action = "singleHostTagResolve"
    c.httpRequest.Body = request
    var returnValue = &model.SingleHostTagResolveResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListRecordOpenapi 查询解析记录Openapi
func (c *Client) ListRecordOpenapi(request *model.ListRecordOpenapiRequest) (*model.ListRecordOpenapiResponse, error) {
    c.httpRequest.Action = "listRecordOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ListRecordOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// BatchHostTagResolve 批量域名解析
func (c *Client) BatchHostTagResolve(request *model.BatchHostTagResolveRequest) (*model.BatchHostTagResolveResponse, error) {
    c.httpRequest.Action = "batchHostTagResolve"
    c.httpRequest.Body = request
    var returnValue = &model.BatchHostTagResolveResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// GetRecordOpenapi 查询单条解析记录Openapi
func (c *Client) GetRecordOpenapi(request *model.GetRecordOpenapiRequest) (*model.GetRecordOpenapiResponse, error) {
    c.httpRequest.Action = "getRecordOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.GetRecordOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ModifyBalanceRateOpenapi 编辑负载均衡Openapi
func (c *Client) ModifyBalanceRateOpenapi(request *model.ModifyBalanceRateOpenapiRequest) (*model.ModifyBalanceRateOpenapiResponse, error) {
    c.httpRequest.Action = "modifyBalanceRateOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ModifyBalanceRateOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ControlBalanceOpenapi 启停负载均衡Openapi
func (c *Client) ControlBalanceOpenapi(request *model.ControlBalanceOpenapiRequest) (*model.ControlBalanceOpenapiResponse, error) {
    c.httpRequest.Action = "controlBalanceOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ControlBalanceOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListBalanceOpenapi 负载均衡概览Openapi
func (c *Client) ListBalanceOpenapi(request *model.ListBalanceOpenapiRequest) (*model.ListBalanceOpenapiResponse, error) {
    c.httpRequest.Action = "listBalanceOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ListBalanceOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// GetDomainByNameForOpenapi 查询单条域名信息V2
func (c *Client) GetDomainByNameForOpenapi(request *model.GetDomainByNameForOpenapiRequest) (*model.GetDomainByNameForOpenapiResponse, error) {
    c.httpRequest.Action = "getDomainByNameForOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.GetDomainByNameForOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// AddOrder 新增订单
func (c *Client) AddOrder(request *model.AddOrderRequest) (*model.AddOrderResponse, error) {
    c.httpRequest.Action = "addOrder"
    c.httpRequest.Body = request
    var returnValue = &model.AddOrderResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// GetDomainByName 查询单条域名信息
func (c *Client) GetDomainByName(request *model.GetDomainByNameRequest) (*model.GetDomainByNameResponse, error) {
    c.httpRequest.Action = "getDomainByName"
    c.httpRequest.Body = request
    var returnValue = &model.GetDomainByNameResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// GetBalanceDetailOpenapi 查询负载均衡详情Openapi
func (c *Client) GetBalanceDetailOpenapi(request *model.GetBalanceDetailOpenapiRequest) (*model.GetBalanceDetailOpenapiResponse, error) {
    c.httpRequest.Action = "getBalanceDetailOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.GetBalanceDetailOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// CreateDomainOpenapi 新增域名V2
func (c *Client) CreateDomainOpenapi(request *model.CreateDomainOpenapiRequest) (*model.CreateDomainOpenapiResponse, error) {
    c.httpRequest.Action = "createDomainOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.CreateDomainOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListNs 查询云解析NS服务器
func (c *Client) ListNs() (*model.ListNsResponse, error) {
    c.httpRequest.Action = "listNs"
    
    var returnValue = &model.ListNsResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// CancelDnsOrder 退订
func (c *Client) CancelDnsOrder(request *model.CancelDnsOrderRequest) (*model.CancelDnsOrderResponse, error) {
    c.httpRequest.Action = "cancelDnsOrder"
    c.httpRequest.Body = request
    var returnValue = &model.CancelDnsOrderResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// DomainStatisticsOpenapi 根据域名查询统计量openapi
func (c *Client) DomainStatisticsOpenapi(request *model.DomainStatisticsOpenapiRequest) (*model.DomainStatisticsOpenapiResponse, error) {
    c.httpRequest.Action = "domainStatisticsOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.DomainStatisticsOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ControlDomainOpenapi 启停域名V2
func (c *Client) ControlDomainOpenapi(request *model.ControlDomainOpenapiRequest) (*model.ControlDomainOpenapiResponse, error) {
    c.httpRequest.Action = "controlDomainOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ControlDomainOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// RenewDnsOrder 续订
func (c *Client) RenewDnsOrder(request *model.RenewDnsOrderRequest) (*model.RenewDnsOrderResponse, error) {
    c.httpRequest.Action = "renewDnsOrder"
    c.httpRequest.Body = request
    var returnValue = &model.RenewDnsOrderResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// CreateFreeDomain 新增免费域名
func (c *Client) CreateFreeDomain(request *model.CreateFreeDomainRequest) (*model.CreateFreeDomainResponse, error) {
    c.httpRequest.Action = "createFreeDomain"
    c.httpRequest.Body = request
    var returnValue = &model.CreateFreeDomainResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// DeleteDomainOpenapi 删除域名V2
func (c *Client) DeleteDomainOpenapi(request *model.DeleteDomainOpenapiRequest) (*model.DeleteDomainOpenapiResponse, error) {
    c.httpRequest.Action = "deleteDomainOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.DeleteDomainOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ModifyOrder 变更订单
func (c *Client) ModifyOrder(request *model.ModifyOrderRequest) (*model.ModifyOrderResponse, error) {
    c.httpRequest.Action = "modifyOrder"
    c.httpRequest.Body = request
    var returnValue = &model.ModifyOrderResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ModifyDomainDesc 修改域名描述
func (c *Client) ModifyDomainDesc(request *model.ModifyDomainDescRequest) (*model.ModifyDomainDescResponse, error) {
    c.httpRequest.Action = "modifyDomainDesc"
    c.httpRequest.Body = request
    var returnValue = &model.ModifyDomainDescResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ModifyDomainDescOpenapi 修改域名描述V2
func (c *Client) ModifyDomainDescOpenapi(request *model.ModifyDomainDescOpenapiRequest) (*model.ModifyDomainDescOpenapiResponse, error) {
    c.httpRequest.Action = "modifyDomainDescOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ModifyDomainDescOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListInstance 获取实例列表
func (c *Client) ListInstance(request *model.ListInstanceRequest) (*model.ListInstanceResponse, error) {
    c.httpRequest.Action = "listInstance"
    c.httpRequest.Body = request
    var returnValue = &model.ListInstanceResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// CreateFreeDomainOpenapi 新增免费域名V2
func (c *Client) CreateFreeDomainOpenapi(request *model.CreateFreeDomainOpenapiRequest) (*model.CreateFreeDomainOpenapiResponse, error) {
    c.httpRequest.Action = "createFreeDomainOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.CreateFreeDomainOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListNsOpenapi 查询云解析NS服务器Openapi
func (c *Client) ListNsOpenapi() (*model.ListNsOpenapiResponse, error) {
    c.httpRequest.Action = "listNsOpenapi"
    
    var returnValue = &model.ListNsOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ModifyDomainOpenapi 修改实例关联的域名V2
func (c *Client) ModifyDomainOpenapi(request *model.ModifyDomainOpenapiRequest) (*model.ModifyDomainOpenapiResponse, error) {
    c.httpRequest.Action = "modifyDomainOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ModifyDomainOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// CreateOrderOpenapi 下单Openapi
func (c *Client) CreateOrderOpenapi(request *model.CreateOrderOpenapiRequest) (*model.CreateOrderOpenapiResponse, error) {
    c.httpRequest.Action = "createOrderOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.CreateOrderOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// GetOperationLog 获取操作日志
func (c *Client) GetOperationLog(request *model.GetOperationLogRequest) (*model.GetOperationLogResponse, error) {
    c.httpRequest.Action = "getOperationLog"
    c.httpRequest.Body = request
    var returnValue = &model.GetOperationLogResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListDomainForOpenapi 显示域名列表V2
func (c *Client) ListDomainForOpenapi(request *model.ListDomainForOpenapiRequest) (*model.ListDomainForOpenapiResponse, error) {
    c.httpRequest.Action = "listDomainForOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ListDomainForOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// CancelOrderForOpenapi 退订Openapi
func (c *Client) CancelOrderForOpenapi(request *model.CancelOrderForOpenapiRequest) (*model.CancelOrderForOpenapiResponse, error) {
    c.httpRequest.Action = "cancelOrderForOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.CancelOrderForOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ModifyDomain 修改实例关联的域名
func (c *Client) ModifyDomain(request *model.ModifyDomainRequest) (*model.ModifyDomainResponse, error) {
    c.httpRequest.Action = "modifyDomain"
    c.httpRequest.Body = request
    var returnValue = &model.ModifyDomainResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// LockDomainOpenapi 域名锁定V2
func (c *Client) LockDomainOpenapi(request *model.LockDomainOpenapiRequest) (*model.LockDomainOpenapiResponse, error) {
    c.httpRequest.Action = "lockDomainOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.LockDomainOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListDomain 显示域名列表
func (c *Client) ListDomain(request *model.ListDomainRequest) (*model.ListDomainResponse, error) {
    c.httpRequest.Action = "listDomain"
    c.httpRequest.Body = request
    var returnValue = &model.ListDomainResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// RenewProductOpenapi 续订产品Openapi
func (c *Client) RenewProductOpenapi(request *model.RenewProductOpenapiRequest) (*model.RenewProductOpenapiResponse, error) {
    c.httpRequest.Action = "renewProductOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.RenewProductOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// GetOperationLogOpenapi 获取操作日志V2
func (c *Client) GetOperationLogOpenapi(request *model.GetOperationLogOpenapiRequest) (*model.GetOperationLogOpenapiResponse, error) {
    c.httpRequest.Action = "getOperationLogOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.GetOperationLogOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// TakeOverRefreshDomain 域名接管刷新
func (c *Client) TakeOverRefreshDomain(request *model.TakeOverRefreshDomainRequest) (*model.TakeOverRefreshDomainResponse, error) {
    c.httpRequest.Action = "takeOverRefreshDomain"
    c.httpRequest.Body = request
    var returnValue = &model.TakeOverRefreshDomainResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListInstanceForOpenapi 获取实例列表openapi
func (c *Client) ListInstanceForOpenapi(request *model.ListInstanceForOpenapiRequest) (*model.ListInstanceForOpenapiResponse, error) {
    c.httpRequest.Action = "listInstanceForOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ListInstanceForOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// TakeOverRefreshDomainOpenapi 域名接管刷新V2
func (c *Client) TakeOverRefreshDomainOpenapi(request *model.TakeOverRefreshDomainOpenapiRequest) (*model.TakeOverRefreshDomainOpenapiResponse, error) {
    c.httpRequest.Action = "takeOverRefreshDomainOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.TakeOverRefreshDomainOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// UpgradeOrderOpenapi 升级套餐Openapi
func (c *Client) UpgradeOrderOpenapi(request *model.UpgradeOrderOpenapiRequest) (*model.UpgradeOrderOpenapiResponse, error) {
    c.httpRequest.Action = "upgradeOrderOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.UpgradeOrderOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// UnlockDomainOpenapi 域名解锁V2
func (c *Client) UnlockDomainOpenapi(request *model.UnlockDomainOpenapiRequest) (*model.UnlockDomainOpenapiResponse, error) {
    c.httpRequest.Action = "unlockDomainOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.UnlockDomainOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ModifyRecord 修改解析记录
func (c *Client) ModifyRecord(request *model.ModifyRecordRequest) (*model.ModifyRecordResponse, error) {
    c.httpRequest.Action = "modifyRecord"
    c.httpRequest.Body = request
    var returnValue = &model.ModifyRecordResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// CreateRecordOpenapi 新增解析记录Openapi
func (c *Client) CreateRecordOpenapi(request *model.CreateRecordOpenapiRequest) (*model.CreateRecordOpenapiResponse, error) {
    c.httpRequest.Action = "createRecordOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.CreateRecordOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// CreateCustomLineOpenapi 创建自定义线路V2
func (c *Client) CreateCustomLineOpenapi(request *model.CreateCustomLineOpenapiRequest) (*model.CreateCustomLineOpenapiResponse, error) {
    c.httpRequest.Action = "createCustomLineOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.CreateCustomLineOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// CreateRecord 新增解析记录
func (c *Client) CreateRecord(request *model.CreateRecordRequest) (*model.CreateRecordResponse, error) {
    c.httpRequest.Action = "createRecord"
    c.httpRequest.Body = request
    var returnValue = &model.CreateRecordResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// RemoveCustomLineOpenapi 删除自定义线路V2
func (c *Client) RemoveCustomLineOpenapi(request *model.RemoveCustomLineOpenapiRequest) (*model.RemoveCustomLineOpenapiResponse, error) {
    c.httpRequest.Action = "removeCustomLineOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.RemoveCustomLineOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListCustomLineOpenapi 查询自定义线路列表V2
func (c *Client) ListCustomLineOpenapi(request *model.ListCustomLineOpenapiRequest) (*model.ListCustomLineOpenapiResponse, error) {
    c.httpRequest.Action = "listCustomLineOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ListCustomLineOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListNewRecordLine 查询解析线路-树状结构
func (c *Client) ListNewRecordLine(request *model.ListNewRecordLineRequest) (*model.ListNewRecordLineResponse, error) {
    c.httpRequest.Action = "listNewRecordLine"
    c.httpRequest.Body = request
    var returnValue = &model.ListNewRecordLineResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// CreateLineGroupOpenapi 创建线路分组V2
func (c *Client) CreateLineGroupOpenapi(request *model.CreateLineGroupOpenapiRequest) (*model.CreateLineGroupOpenapiResponse, error) {
    c.httpRequest.Action = "createLineGroupOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.CreateLineGroupOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListRecord 查询解析记录
func (c *Client) ListRecord(request *model.ListRecordRequest) (*model.ListRecordResponse, error) {
    c.httpRequest.Action = "listRecord"
    c.httpRequest.Body = request
    var returnValue = &model.ListRecordResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// RemoveLineGroupOpenapi 删除线路分组列表V2
func (c *Client) RemoveLineGroupOpenapi(request *model.RemoveLineGroupOpenapiRequest) (*model.RemoveLineGroupOpenapiResponse, error) {
    c.httpRequest.Action = "removeLineGroupOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.RemoveLineGroupOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// GetRecord 查询单条解析记录
func (c *Client) GetRecord(request *model.GetRecordRequest) (*model.GetRecordResponse, error) {
    c.httpRequest.Action = "getRecord"
    c.httpRequest.Body = request
    var returnValue = &model.GetRecordResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListLineGroupOpenapi 查询线路分组列表V2
func (c *Client) ListLineGroupOpenapi(request *model.ListLineGroupOpenapiRequest) (*model.ListLineGroupOpenapiResponse, error) {
    c.httpRequest.Action = "listLineGroupOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ListLineGroupOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ModifyLineGroupOpenapi 修改线路分组V2
func (c *Client) ModifyLineGroupOpenapi(request *model.ModifyLineGroupOpenapiRequest) (*model.ModifyLineGroupOpenapiResponse, error) {
    c.httpRequest.Action = "modifyLineGroupOpenapi"
    c.httpRequest.Body = request
    var returnValue = &model.ModifyLineGroupOpenapiResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ModifyBalanceRate 编辑负载均衡
func (c *Client) ModifyBalanceRate(request *model.ModifyBalanceRateRequest) (*model.ModifyBalanceRateResponse, error) {
    c.httpRequest.Action = "modifyBalanceRate"
    c.httpRequest.Body = request
    var returnValue = &model.ModifyBalanceRateResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ControlBalance 启停负载均衡
func (c *Client) ControlBalance(request *model.ControlBalanceRequest) (*model.ControlBalanceResponse, error) {
    c.httpRequest.Action = "controlBalance"
    c.httpRequest.Body = request
    var returnValue = &model.ControlBalanceResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// GetBalanceDetail 查询负载均衡详情
func (c *Client) GetBalanceDetail(request *model.GetBalanceDetailRequest) (*model.GetBalanceDetailResponse, error) {
    c.httpRequest.Action = "getBalanceDetail"
    c.httpRequest.Body = request
    var returnValue = &model.GetBalanceDetailResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// RemoveLineGroup 删除线路分组列表
func (c *Client) RemoveLineGroup(request *model.RemoveLineGroupRequest) (*model.RemoveLineGroupResponse, error) {
    c.httpRequest.Action = "removeLineGroup"
    c.httpRequest.Body = request
    var returnValue = &model.RemoveLineGroupResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListBalance 负载均衡概览
func (c *Client) ListBalance(request *model.ListBalanceRequest) (*model.ListBalanceResponse, error) {
    c.httpRequest.Action = "listBalance"
    c.httpRequest.Body = request
    var returnValue = &model.ListBalanceResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// CreateLineGroup 创建线路分组
func (c *Client) CreateLineGroup(request *model.CreateLineGroupRequest) (*model.CreateLineGroupResponse, error) {
    c.httpRequest.Action = "createLineGroup"
    c.httpRequest.Body = request
    var returnValue = &model.CreateLineGroupResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ListLineGroup 查询线路分组列表
func (c *Client) ListLineGroup(request *model.ListLineGroupRequest) (*model.ListLineGroupResponse, error) {
    c.httpRequest.Action = "listLineGroup"
    c.httpRequest.Body = request
    var returnValue = &model.ListLineGroupResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ModifyLineGroup 修改线路分组
func (c *Client) ModifyLineGroup(request *model.ModifyLineGroupRequest) (*model.ModifyLineGroupResponse, error) {
    c.httpRequest.Action = "modifyLineGroup"
    c.httpRequest.Body = request
    var returnValue = &model.ModifyLineGroupResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
// ControlDomain 启停域名
func (c *Client) ControlDomain(request *model.ControlDomainRequest) (*model.ControlDomainResponse, error) {
    c.httpRequest.Action = "controlDomain"
    c.httpRequest.Body = request
    var returnValue = &model.ControlDomainResponse{}
    if _, err := c.APIClient.Excute(c.httpRequest, c.config, returnValue); err != nil {
        return nil, err
    } else {
        return returnValue, nil
    }
}
