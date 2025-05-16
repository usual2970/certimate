package edgio_api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Edgio/edgio-api/applications/v7/dtos"
	"github.com/go-resty/resty/v2"
)

// AccessTokenResponse represents the response from the token endpoint.
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// TokenCache represents a cached token. The token is stored along
// with its expiry time. Because different endpoints require different
// scopes, we store the token with the scope as the key, so that we
// can fetch the token from the cache based on the scope.
type TokenCache struct {
	AccessToken string
	Expiry      time.Time
}

type EdgioClient struct {
	client       *resty.Client
	clientID     string
	clientSecret string
	tokenURL     string
	apiURL       string
	tokenCache   map[string]TokenCache
}

func NewEdgioClient(clientID, clientSecret, tokenURL, apiURL string) *EdgioClient {
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second)

	if tokenURL == "" {
		tokenURL = "https://id.edgio.app/connect/token"
	}

	if apiURL == "" {
		apiURL = "https://edgioapis.com"
	}

	return &EdgioClient{
		client:       client,
		clientID:     clientID,
		clientSecret: clientSecret,
		tokenURL:     tokenURL,
		apiURL:       apiURL,
		tokenCache:   make(map[string]TokenCache),
	}
}

func (c *EdgioClient) getToken(scope string) (string, error) {
	if cachedToken, exists := c.tokenCache[scope]; exists && time.Now().Before(cachedToken.Expiry) {
		return cachedToken.AccessToken, nil
	}

	var tokenResp AccessTokenResponse
	resp, err := c.client.R().
		SetFormData(map[string]string{
			"client_id":     c.clientID,
			"client_secret": c.clientSecret,
			"grant_type":    "client_credentials",
			"scope":         scope,
		}).
		SetResult(&tokenResp).
		Post(c.tokenURL)
	if err != nil {
		return "", fmt.Errorf("failed to request token: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("unexpected status code for getToken: %d", resp.StatusCode())
	}

	c.tokenCache[scope] = TokenCache{
		AccessToken: tokenResp.AccessToken,
		Expiry:      time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}

	return tokenResp.AccessToken, nil
}

func (c *EdgioClient) GetProperty(ctx context.Context, propertyID string) (*dtos.Property, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/properties/%s", c.apiURL, propertyID)

	var property dtos.Property
	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetResult(&property).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status code for getSpecificProperty: %d, %s", resp.StatusCode(), resp.Request.URL)
	}

	return &property, nil
}

func (c *EdgioClient) GetProperties(page int, pageSize int, organizationID string) (*dtos.Properties, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/properties", c.apiURL)

	var propertiesResp dtos.Properties
	resp, err := c.client.R().
		SetAuthToken(token).
		SetQueryParams(map[string]string{
			"page":            fmt.Sprintf("%d", page),
			"page_size":       fmt.Sprintf("%d", pageSize),
			"organization_id": organizationID,
		}).
		SetResult(&propertiesResp).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status code for getProperties: %d, %s", resp.StatusCode(), resp.String())
	}

	return &propertiesResp, nil
}

func (c *EdgioClient) CreateProperty(ctx context.Context, organizationID, slug string) (*dtos.Property, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/properties", c.apiURL)

	var createdProperty dtos.Property
	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"organization_id": organizationID,
			"slug":            slug,
		}).
		SetResult(&createdProperty).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status code for createProperty: %d, response: %s", resp.StatusCode(), resp.String())
	}

	return &createdProperty, nil
}

func (c *EdgioClient) DeleteProperty(propertyID string) error {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/properties/%s", c.apiURL, propertyID)

	resp, err := c.client.R().
		SetAuthToken(token).
		Delete(url)
	if err != nil {
		return fmt.Errorf("error sending DELETE request: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("error deleting property: status code %d", resp.StatusCode())
	}

	return nil
}

func (c *EdgioClient) UpdateProperty(ctx context.Context, propertyID string, slug string) (*dtos.Property, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/properties/%s", c.apiURL, propertyID)

	requestBody := map[string]interface{}{
		"slug": slug,
	}

	var updatedProperty dtos.Property
	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(requestBody).
		SetResult(&updatedProperty).
		Patch(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status code for updateProperty: %d", resp.StatusCode())
	}

	return &updatedProperty, nil
}

func (c *EdgioClient) GetEnvironments(page, pageSize int, propertyID string) (*dtos.EnvironmentsResponse, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/environments", c.apiURL)

	resp, err := c.client.R().
		SetAuthToken(token).
		SetQueryParams(map[string]string{
			"page":        fmt.Sprintf("%d", page),
			"page_size":   fmt.Sprintf("%d", pageSize),
			"property_id": propertyID,
		}).
		SetResult(&dtos.EnvironmentsResponse{}).
		Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.String())
	}

	return resp.Result().(*dtos.EnvironmentsResponse), nil
}

func (c *EdgioClient) GetEnvironment(environmentID string) (*dtos.Environment, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/environments/%s", c.apiURL, environmentID)

	resp, err := c.client.R().
		SetPathParams(map[string]string{
			"environment_id": environmentID,
		}).
		SetAuthToken(token).
		SetResult(&dtos.Environment{}).
		Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.String())
	}

	return resp.Result().(*dtos.Environment), nil
}

func (c *EdgioClient) CreateEnvironment(propertyID, name string, onlyMaintainersCanDeploy, httpRequestLogging bool) (*dtos.Environment, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/environments", c.apiURL)

	body := map[string]interface{}{
		"property_id":                 propertyID,
		"name":                        name,
		"only_maintainers_can_deploy": onlyMaintainersCanDeploy,
		"http_request_logging":        httpRequestLogging,
	}

	resp, err := c.client.R().
		SetBody(body).
		SetAuthToken(token).
		SetResult(&dtos.Environment{}).
		Post(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.String())
	}

	return resp.Result().(*dtos.Environment), nil
}

func (c *EdgioClient) UpdateEnvironment(environmentID, name string, onlyMaintainersCanDeploy, httpRequestLogging, preserveCache bool) (*dtos.Environment, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/environments/%s", c.apiURL, environmentID)

	body := map[string]interface{}{
		"name": name,
		// as can_members_deploy is depricated, but update api is not
		// we need to use it to map onlyMaintainersCanDeploy
		"only_maintainers_can_deploy": onlyMaintainersCanDeploy,
		"http_request_logging":        httpRequestLogging,
		"preserve_cache":              preserveCache,
	}

	resp, err := c.client.R().
		SetPathParams(map[string]string{
			"environment_id": environmentID,
		}).
		SetBody(body).
		SetAuthToken(token).
		SetResult(&dtos.Environment{}).
		Patch(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.String())
	}

	return resp.Result().(*dtos.Environment), nil
}

func (c *EdgioClient) DeleteEnvironment(environmentID string) error {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/environments/%s", c.apiURL, environmentID)

	resp, err := c.client.R().
		SetPathParams(map[string]string{
			"environment_id": environmentID,
		}).
		SetAuthToken(token).
		SetResult(&dtos.Environment{}).
		Delete(url)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("error response: %s", resp.String())
	}

	return nil
}

func (c *EdgioClient) GetTlsCert(tlsCertId string) (*dtos.TLSCertResponse, error) {
	token, err := c.getToken("app.config")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/config/v0.1/tls-certs/%s", c.apiURL, tlsCertId)

	var tlsCertResponse dtos.TLSCertResponse
	resp, err := c.client.R().
		SetAuthToken(token).
		SetResult(&tlsCertResponse).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("error response: %s", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.String())
	}

	return &tlsCertResponse, nil
}

func (c *EdgioClient) UploadTlsCert(req dtos.UploadTlsCertRequest) (*dtos.TLSCertResponse, error) {
	token, err := c.getToken("app.config")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/config/v0.1/tls-certs", c.apiURL)
	response := &dtos.TLSCertResponse{}

	resp, err := c.client.R().
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(response).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("failed to upload TLS certificate: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API responded with error: %s", resp.String())
	}

	return response, nil
}

func (c *EdgioClient) GenerateTlsCert(environmentId string) (*dtos.TLSCertResponse, error) {
	token, err := c.getToken("app.config")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/config/v0.1/tls-certs/generate", c.apiURL)
	request := map[string]interface{}{
		"environment_id": environmentId,
	}
	response := &dtos.TLSCertResponse{}

	resp, err := c.client.R().
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(response).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("failed to upload TLS certificate: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API responded with error: %s", resp.String())
	}

	return response, nil
}

func (c *EdgioClient) GetTlsCerts(page int, pageSize int, environmentID string) (*dtos.TLSCertSResponse, error) {
	token, err := c.getToken("app.config")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/config/v0.1/tls-certs", c.apiURL)

	var tlsCertsResponse dtos.TLSCertSResponse
	resp, err := c.client.R().
		SetAuthToken(token).
		SetQueryParams(map[string]string{
			"page":           fmt.Sprintf("%d", page),
			"page_size":      fmt.Sprintf("%d", pageSize),
			"environment_id": environmentID,
		}).
		SetResult(&tlsCertsResponse).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status code for getTlsCerts: %d", resp.StatusCode())
	}

	return &tlsCertsResponse, nil
}

func (c *EdgioClient) UploadCdnConfiguration(config *dtos.CDNConfiguration) (*dtos.CDNConfiguration, error) {
	fmt.Println("------------------------------------------------------------------------- uploading")

	token, err := c.getToken("app.config")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/config/v0.1/configs", c.apiURL)
	var response dtos.CDNConfiguration

	// Convert config to json
	jsonBody, _ := json.MarshalIndent(config, "", "    ")
	jsonString := string(jsonBody)
	fmt.Println("------------------------- config report code: ", config.Hostnames[0].ReportCode == nil)
	fmt.Println("------------------------- config report code value: ", config.Hostnames[0].ReportCode)
	fmt.Println("----------------------------------- jsonBody: ", jsonString)

	resp, err := c.client.R().
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(config).
		SetResult(&response).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("failed to upload CDN configuration: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status code for uploadCdnConfiguration: %d, %s", resp.StatusCode(), resp.String())
	}

	return &response, nil
}

func (c *EdgioClient) GetCDNConfiguration(configID string) (*dtos.CDNConfiguration, error) {
	fmt.Println("------------------------------------------------------------------------- reading config")

	token, err := c.getToken("app.config")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("https://edgioapis.com/config/v0.1/configs/%s", configID)
	var response dtos.CDNConfiguration

	resp, err := c.client.R().
		SetAuthToken(token).
		SetResult(&response).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get CDN configuration: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status code for GetCDNConfiguration: %d", resp.StatusCode())
	}

	return &response, nil
}
