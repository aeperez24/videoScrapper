package cuevana

import (
	"aeperez24/animewatcher/service"
	"fmt"
	"io"
	"net/url"
	"strings"
)

const proxies_url string = "https://www.proxy-list.download/api/v2/get?l=en&t=https"

type DownloaderService struct {
	ScrapService         ScrapperService
	GetSender            service.GetSender
	proxies              []string
	getProxies           func() []string
	getClientWithProxies func(*[]string) httpPostClient
}

func NewDownloaderService(httpWraper service.GetSender) DownloaderService {
	return DownloaderService{
		ScrapService:         ScrapperService{},
		GetSender:            httpWraper,
		getProxies:           getProxies,
		getClientWithProxies: getHttpClientWithProxy,
	}

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
	cl := ds.getClientWithProxies(&ds.proxies)
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
