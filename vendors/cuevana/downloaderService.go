package cuevana

import (
	"aeperez24/animewatcher/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const proxies_url string = "https://www.proxy-list.download/api/v2/get?l=en&t=https"

type DownloaderService struct {
	ScrapService ScrapperService
	GetSender    service.GetSender
	proxies      []string
}

func (ds DownloaderService) GetSortedEpisodesAvaliable(serieLink string) ([]string, error) {

	linksList, err := ds.GetSortedEpisodesLinks(serieLink)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, link := range linksList {
		result = append(result, ds.ScrapService.getEpisodeName(link))
	}
	return result, err
}

func (ds DownloaderService) GetSortedEpisodesLinks(serieLink string) ([]string, error) {
	resp, err := ds.GetSender.Get(serieLink)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ds.ScrapService.getEpisodesList(resp.Body)
}
func (ds DownloaderService) DownloadEpisodeFromLink(serieLink string, episodeNumber string) (io.Reader, string, error) {
	episodesLinks, _ := ds.GetSortedEpisodesLinks(serieLink)
	for _, episodeLink := range episodesLinks {
		if strings.Contains(episodeLink, episodeNumber) {
			episodePage, err := ds.GetSender.Get(episodeLink)
			if err != nil {
				return nil, "", err
			}
			defer episodePage.Body.Close()
			fichierLink, _ := ds.ScrapService.get1fichierLink(episodePage.Body)
			fichierPage, _ := ds.GetSender.Get(fichierLink)
			defer fichierPage.Body.Close()

			params, err := ds.ScrapService.getParammsForFichierDownload(fichierPage.Body)
			if err != nil {
				return nil, "", err
			}
			resultVideo, err := ds.downloadFromFichier(params.postUrl, params.adz)
			return resultVideo, "mp4", err

		}
	}
	return nil, "", fmt.Errorf("episode %v not found for link %v", episodeNumber, serieLink)
}

func (ds DownloaderService) downloadFromFichier(fichierLink string, adz string) (io.Reader, error) {
	cl := ds.getHttpClientWithProxy()
	resp, err := cl.PostForm(fichierLink, url.Values{"adz": {adz}, "dl_no_ssl": {"on"}, "dlinline": {"on"}})
	if err != nil {
		return nil, err
	}
	downloadLink, err := ds.ScrapService.getDownloadLink(resp.Body)

	if err != nil {
		return nil, err
	}

	res, err := ds.GetSender.Get(downloadLink)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func (ds *DownloaderService) getHttpClientWithProxy() *http.Client {
	if len(ds.proxies) == 0 {
		ds.proxies = getProxies()
	}
	proxy := ds.proxies[0]
	ds.proxies = ds.proxies[1:]
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
