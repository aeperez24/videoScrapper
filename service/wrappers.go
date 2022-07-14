package service

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type HttpWrapper interface {
	Get(url string) (*http.Response, error)
	PostForm(url string, formValues url.Values) (*http.Response, error)
	Post(urlx string, contentType string, reader io.Reader) (*http.Response, error)
	Request(urlx string, method string, reader io.Reader) (*http.Response, error)
	RequestWithHeaders(urlx string, method string, reader io.Reader, headers map[string]string) (*http.Response, error)
}
type HttpWrapperImpl struct {
}

type FileSystemManager interface {
	Save(string, string, io.Reader) error
	Read(path string, fileName string) ([]byte, error)
}

func (wrapper FileSystemManagerWrapper) Save(filepath string, fileName string, reader io.Reader) error {
	os.MkdirAll(filepath, os.ModePerm)
	out, err := os.Create(filepath + "/" + fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, reader)
	log.Printf("file %v saved on %v", fileName, filepath)
	return nil
}

func (wrapper FileSystemManagerWrapper) Read(filepath string, fileName string) ([]byte, error) {
	return os.ReadFile(filepath + "/" + fileName)
}

type FileSystemManagerWrapper struct {
}

func (wrapper HttpWrapperImpl) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func (wrapper HttpWrapperImpl) PostForm(urlx string, formValues url.Values) (*http.Response, error) {
	log.Println("posting to:" + urlx)
	return http.PostForm(urlx, formValues)
}

func (wrapper HttpWrapperImpl) Post(urlx string, contentType string, reader io.Reader) (*http.Response, error) {
	return http.Post(urlx, contentType, reader)
}

func (wrapper HttpWrapperImpl) Request(urlx string, method string, reader io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(method, urlx, reader)
	cl := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	return cl.Do(req)
}

func (wrapper HttpWrapperImpl) RequestWithHeaders(urlx string, method string, reader io.Reader, headers map[string]string) (*http.Response, error) {
	req, _ := http.NewRequest(method, urlx, reader)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cl := http.Client{
		Transport: tr,
	}
	for key, element := range headers {
		req.Header.Add(key, element)
	}
	return cl.Do(req)
}
