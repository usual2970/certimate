package http

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
)

type Options struct {
	Timeout time.Duration
}

type Option func(o *Options)

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

func Req(url string, method string, body io.Reader, head map[string]string, opts ...Option) ([]byte, error) {
	reader, err := Req2GetReader(url, method, body, head, opts...)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

func Req2GetReader(url string, method string, body io.Reader, head map[string]string, opts ...Option) (io.ReadCloser, error) {
	req := BuildReq(url, method, body, head)

	return ToRequest(req, opts...)

}

func BuildReq(url string, method string, body io.Reader, head map[string]string) *http.Request {

	// Create an http.Request instance
	req, _ := http.NewRequest(method, url, body)
	for k, v := range head {
		req.Header.Set(k, v)
	}

	return req
}

func ToRequest(req *http.Request, opts ...Option) (io.ReadCloser, error) {
	options := &Options{
		Timeout: 30000 * time.Millisecond,
	}

	for _, opt := range opts {
		opt(options)
	}

	client := httpclient.NewClient(httpclient.WithHTTPTimeout(options.Timeout))
	// Call the `Do` method, which has a similar interface to the `http.Do` method
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is not 200: %d", res.StatusCode)
	}

	return res.Body, nil
}
