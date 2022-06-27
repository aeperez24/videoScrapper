package cuevana

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
)

func getHttpClientWithProxy(proxies *[]string) httpPostClient {
	nproxies := *proxies
	if len(*proxies) == 0 {
		nproxies = getProxies()
	}
	randNumber := rand.Intn(len(nproxies))
	proxy := (nproxies)[randNumber]
	nproxies = removeElement(nproxies, randNumber)
	proxies = &(nproxies)
	proxyUrl, _ := url.Parse(proxy)
	return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
}

func getProxies() []string {
	resp, _ := http.Get(proxies_url)
	presponse := proxyResponse{}
	json.NewDecoder(resp.Body).Decode(&presponse)
	result := make([]string, 0)
	for _, proxy := range presponse.LISTA {
		result = append(result, fmt.Sprintf("http://%s:%s", proxy.IP, proxy.PORT))
	}
	return result
}

type proxy struct {
	IP   string
	PORT string
}

type proxyResponse struct {
	LISTA []proxy
}

type httpPostClient interface {
	PostForm(url string, data url.Values) (resp *http.Response, err error)
	Get(url string) (*http.Response, error)
}

func removeElement(in []string, pos int) []string {
	return append(in[0:pos], in[pos:]...)

}
