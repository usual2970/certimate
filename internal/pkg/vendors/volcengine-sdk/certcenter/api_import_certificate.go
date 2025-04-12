package certcenter

import (
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"github.com/volcengine/volcengine-go-sdk/volcengine/request"
	"github.com/volcengine/volcengine-go-sdk/volcengine/response"
	"github.com/volcengine/volcengine-go-sdk/volcengine/volcengineutil"
)

const opImportCertificateCommon = "ImportCertificate"

func (c *CertCenter) ImportCertificateCommonRequest(input *map[string]interface{}) (req *request.Request, output *map[string]interface{}) {
	op := &request.Operation{
		Name:       opImportCertificateCommon,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &map[string]interface{}{}
	}

	output = &map[string]interface{}{}
	req = c.newRequest(op, input, output)

	req.HTTPRequest.Header.Set("Content-Type", "application/json; charset=utf-8")

	return
}

func (c *CertCenter) ImportCertificateCommon(input *map[string]interface{}) (*map[string]interface{}, error) {
	req, out := c.ImportCertificateCommonRequest(input)
	return out, req.Send()
}

func (c *CertCenter) ImportCertificateCommonWithContext(ctx volcengine.Context, input *map[string]interface{}, opts ...request.Option) (*map[string]interface{}, error) {
	req, out := c.ImportCertificateCommonRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

const opImportCertificate = "ImportCertificate"

func (c *CertCenter) ImportCertificateRequest(input *ImportCertificateInput) (req *request.Request, output *ImportCertificateOutput) {
	op := &request.Operation{
		Name:       opImportCertificate,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &ImportCertificateInput{}
	}

	output = &ImportCertificateOutput{}
	req = c.newRequest(op, input, output)

	req.HTTPRequest.Header.Set("Content-Type", "application/json; charset=utf-8")

	return
}

func (c *CertCenter) ImportCertificate(input *ImportCertificateInput) (*ImportCertificateOutput, error) {
	req, out := c.ImportCertificateRequest(input)
	return out, req.Send()
}

func (c *CertCenter) ImportCertificateWithContext(ctx volcengine.Context, input *ImportCertificateInput, opts ...request.Option) (*ImportCertificateOutput, error) {
	req, out := c.ImportCertificateRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type ImportCertificateInput struct {
	_ struct{} `type:"structure" json:",omitempty"`

	Tag *string `type:"string" json:",omitempty"`

	ProjectName *string `type:"string" json:",omitempty"`

	Repeatable *bool `type:"boolean" json:",omitempty"`

	NoVerifyAndFixChain *bool `type:"boolean" json:",omitempty"`

	CertificateInfo *ImportCertificateInputCertificateInfo `type:"structure" json:",omitempty"`

	Tags *[]ImportCertificateInputTag `type:"list" json:",omitempty"`
}

func (s ImportCertificateInput) String() string {
	return volcengineutil.Prettify(s)
}

func (s *ImportCertificateInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "ImportCertificateInput"}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

type ImportCertificateInputCertificateInfo struct {
	CertificateChain *string `type:"string" json:",omitempty"`

	PrivateKey *string `type:"string" json:",omitempty"`
}

type ImportCertificateInputTag struct {
	Key *string `type:"string" json:",omitempty" required:"true"`

	Value *string `type:"string" json:",omitempty" required:"true"`
}

type ImportCertificateOutput struct {
	_ struct{} `type:"structure" json:",omitempty"`

	Metadata *response.ResponseMetadata

	InstanceId *string `type:"string" json:",omitempty"`

	RepeatId *string `type:"string" json:",omitempty"`
}

func (s ImportCertificateOutput) String() string {
	return volcengineutil.Prettify(s)
}
