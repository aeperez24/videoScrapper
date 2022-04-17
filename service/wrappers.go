package service

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type GetSender interface {
	Get(url string) (*http.Response, error)
	PostForm(url string, formValues url.Values) (*http.Response, error)
	Post(urlx string, contentType string, reader io.Reader) (*http.Response, error)
	Request(urlx string, method string, reader io.Reader) (*http.Response, error)
	RequestWithHeaders(urlx string, method string, reader io.Reader, headers map[string]string) (*http.Response, error)
}
type GetWrapper struct {
}

type FileSystemSaver interface {
	Save(string, string, io.Reader) error
}

type FileSystemSaverWrapper struct {
}

func (wrapper GetWrapper) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func (wrapper GetWrapper) PostForm(urlx string, formValues url.Values) (*http.Response, error) {
	log.Println("posting to:" + urlx)
	return http.PostForm(urlx, formValues)
}
func (wrapper FileSystemSaverWrapper) Save(filepath string, fileName string, reader io.Reader) error {
	os.MkdirAll(filepath, os.ModePerm)
	out, err := os.Create(filepath + "/" + fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, reader)
	log.Println("file saved on " + filepath)
	return nil
}

func (wrapper GetWrapper) Post(urlx string, contentType string, reader io.Reader) (*http.Response, error) {
	return http.Post(urlx, contentType, reader)
}

func (wrapper GetWrapper) Request(urlx string, method string, reader io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(method, urlx, reader)
	cl := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	return cl.Do(req)
}

func (wrapper GetWrapper) RequestWithHeaders(urlx string, method string, reader io.Reader, headers map[string]string) (*http.Response, error) {
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
