package service

import (
	"io"
	"net/http"
	"net/url"

	"github.com/stretchr/testify/mock"
)

type HttpWrapperMock struct {
	mock.Mock
}

func (wrapper *HttpWrapperMock) Get(url string) (*http.Response, error) {
	args := wrapper.Called(url)
	return buildResponse(args)
}

func (wrapper *HttpWrapperMock) PostForm(url string, formValues url.Values) (*http.Response, error) {
	args := wrapper.Called(url, formValues)
	return buildResponse(args)
}

func (wrapper *HttpWrapperMock) Post(url string, contentType string, reader io.Reader) (*http.Response, error) {
	args := wrapper.Called(url, contentType, reader)
	return buildResponse(args)
}

func (wrapper *HttpWrapperMock) Request(url string, method string, reader io.Reader) (*http.Response, error) {
	args := wrapper.Called(url, method, reader)
	return buildResponse(args)
}

func (wrapper *HttpWrapperMock) RequestWithHeaders(url string, method string, reader io.Reader, headers map[string]string) (*http.Response, error) {
	args := wrapper.Called(url, method, reader, headers)
	return buildResponse(args)
}

func buildResponse(args mock.Arguments) (*http.Response, error) {
	var err error
	err = nil
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return args.Get(0).(*http.Response), err
}
