package http

import (
	"errors"
	"time"

	"myself_framwork/utils/constants"

	"github.com/go-resty/resty/v2"
)

// Env ..
type Env struct {
	DebugClient bool `env:"DEBUG_CLIENT" default:"true"`
	Timeout     int  `env:"TIMEOUT" default:"60s"`
	RetryBad    int  `env:"RETRY_BAD" default:"1"`
}

var (
	httpEnv Env
)

func createNewRequest() *resty.Request {
	client := resty.New().
		SetRetryCount(httpEnv.RetryBad).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		}).
		SetTimeout(time.Duration(httpEnv.Timeout) * time.Second).
		EnableTrace()

	req := client.R()
	return req
}

// HTTPGet func
func HTTPGet(url string, header map[string]string) ([]byte, resty.TraceInfo, error) {
	req := createNewRequest().SetHeaders(header)
	resp, err := req.Execute(constants.HttpMethodGet, url)
	trace := resp.Request.TraceInfo()

	if err != nil {
		return nil, trace, err
	}
	return resp.Body(), trace, nil
}

// HTTPPost func
func HTTPPost(url string, jsondata interface{}, header map[string]string) ([]byte, resty.TraceInfo, error) {
	req := createNewRequest().SetHeaders(header).SetBody(jsondata)
	resp, err := req.Execute(constants.HttpMethodPost, url)
	trace := resp.Request.TraceInfo()

	if err != nil {
		return nil, trace, err
	}
	return resp.Body(), trace, nil
}

// HTTPPutWithHeader func
func HttpPut(url string, jsondata interface{}, header map[string]string) ([]byte, resty.TraceInfo, error) {
	req := createNewRequest().SetHeaders(header).SetBody(jsondata)
	resp, err := req.Execute(constants.HttpMethodPut, url)
	trace := resp.Request.TraceInfo()

	if err != nil {
		return nil, trace, err
	}
	return resp.Body(), trace, nil
}

// HTTPDeleteWithHeader func
func HttpDelete(url string, jsondata interface{}, header map[string]string) ([]byte, resty.TraceInfo, error) {
	req := createNewRequest().SetHeaders(header).SetBody(jsondata)
	resp, err := req.Execute(constants.HttpMethodPost, url)
	trace := resp.Request.TraceInfo()

	if err != nil {
		return nil, trace, err
	}
	return resp.Body(), trace, nil
}

// SendHttpRequest ..
func SendHttpRequest(method string, url string, header map[string]string, body interface{}) ([]byte, resty.TraceInfo, error) {
	var data []byte
	var err error
	var trace resty.TraceInfo
	switch method {
	case constants.HttpMethodGet:
		data, trace, err = HTTPGet(url, header)
	case constants.HttpMethodPost:
		data, trace, err = HTTPPost(url, body, header)
	case constants.HttpMethodPut:
		data, trace, err = HttpPut(url, body, header)
	case constants.HttpMethodDelete:
		data, trace, err = HttpDelete(url, body, header)
	}
	return data, trace, err
}
