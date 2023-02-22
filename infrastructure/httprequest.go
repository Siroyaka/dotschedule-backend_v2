package infrastructure

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

type HTTPRequest struct {
	response       HTTPResponse
	requestTimeout int
}

const (
	defaultHTTPRequestTimeout = 30
)

func (hr HTTPRequest) timeout() time.Duration {
	return time.Second * time.Duration(hr.requestTimeout)
}

func NewHTTPRequest() abstruct.HTTPRequest {
	return &HTTPRequest{
		requestTimeout: defaultHTTPRequestTimeout,
		response:       HTTPResponse{},
	}
}

func (hr HTTPRequest) Get(url string) (abstruct.HTTPResponse, utilerror.IError) {
	client := &http.Client{
		Timeout: hr.timeout(),
	}
	res, err := client.Get(url)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil, utilerror.New(err.Error(), utilerror.ERR_HTTP_REQUEST_TIMEOUT)
	} else if err != nil {
		return nil, utilerror.New(err.Error(), utilerror.ERR_HTTP_REQUEST_ERROR)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, utilerror.New(err.Error(), utilerror.ERR_HTTP_REQUEST_ERROR)
	}
	defer res.Body.Close()
	return HTTPResponse{
		status:     res.Status,
		statusCode: res.StatusCode,
		body:       body,
	}, nil
}

func (hr HTTPRequest) Post(param abstruct.HTTPPostParams) (abstruct.HTTPResponse, utilerror.IError) {
	client := &http.Client{
		Timeout: hr.timeout(),
	}
	res, err := client.Post(param.Url(), param.ContentType(), bytes.NewBufferString(param.Content()))
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil, utilerror.New(err.Error(), utilerror.ERR_HTTP_REQUEST_TIMEOUT)
	} else if err != nil {
		return nil, utilerror.New(err.Error(), utilerror.ERR_HTTP_REQUEST_ERROR)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, utilerror.New(err.Error(), utilerror.ERR_HTTP_REQUEST_ERROR)
	}
	defer res.Body.Close()
	return HTTPResponse{
		status:     res.Status,
		statusCode: res.StatusCode,
		body:       body,
	}, nil
}

func (ht *HTTPRequest) SetTimeout(sec int) {
	ht.requestTimeout = sec
}

type HTTPResponse struct {
	status     string
	statusCode int
	body       []byte
}

func (hr HTTPResponse) Status() string {
	return hr.status
}

func (hr HTTPResponse) StatusCode() int {
	return hr.statusCode
}

func (hr HTTPResponse) Body() string {
	return string(hr.body)
}
