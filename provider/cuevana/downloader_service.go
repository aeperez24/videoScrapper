package cuevana

import (
	"aeperez24/animewatcher/service"
	"fmt"
	"io"
	"net/url"
	"strings"
)

const proxiesUrl string = "https://www.proxy-list.download/api/v2/get?l=en&t=https"

type DownloaderService struct {
	ScrapService         ScrapperService
	HttpWrapper          service.HttpWrapper
	getProxies           func() []string
	getClientWithProxies func([]string) (httpPostClient, string)
	usedProxies          map[string]bool
}

func NewDownloaderService(httpWraper service.HttpWrapper) DownloaderService {
	return DownloaderService{
		ScrapService:         ScrapperService{},
		HttpWrapper:          httpWraper,
		getProxies:           getProxies,
		getClientWithProxies: getHttpClientWithProxy,
		usedProxies:          make(map[string]bool),
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
	resp, err := ds.HttpWrapper.Get(serieLink)
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
			episodePage, err := ds.HttpWrapper.Get(episodeLink)
			if err != nil {
				return nil, "", err
			}

			// We should remove defer functions within loops
			// https://stackoverflow.com/questions/45617758/proper-way-to-release-resources-with-defer-in-a-loop
			defer episodePage.Body.Close()
			fichierLink, _ := ds.ScrapService.get1fichierLink(episodePage.Body)
			fichierPage, _ := ds.HttpWrapper.Get(fichierLink)
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

	proxies := getProxies()

	filtered := ds.filterProxies(proxies)
	cl, proxy := ds.getClientWithProxies(filtered)
	ds.usedProxies[proxy] = true
	resp, err := cl.PostForm(fichierLink, url.Values{"adz": {adz}, "dl_no_ssl": {"on"}, "dlinline": {"on"}})
	if err != nil {
		return nil, err
	}
	downloadLink, err := ds.ScrapService.getDownloadLink(resp.Body)

	if err != nil {
		return nil, err
	}

	res, err := ds.HttpWrapper.Get(downloadLink)
	if err != nil {
		return nil, err
	}
	return res.Body, nil

}

func (ds DownloaderService) filterProxies(proxies []string) []string {
	filtered := make([]string, 0)
	for _, proxy := range proxies {
		if !ds.usedProxies[proxy] {
			filtered = append(filtered, proxy)
		}
	}
	return filtered
}
