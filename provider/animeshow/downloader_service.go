package animeshow

import (
	"aeperez24/videoScrapper/service"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
)

type DowloaderService struct {
	ScrapService     ScrapperService
	HttpWrapper      service.HttpWrapper
	AppConfiguration service.AppConfiguration
}

func (ds DowloaderService) DownloadEpisodeFromLink(serieLink string, episodeNumber string) (io.Reader, string, error) {
	episodes, err := ds.getEpisodesAvaliable(serieLink)

	if err != nil {
		return nil, "", err
	}

	episodeLink, err := ds.getLinkForEpisodeNumber(episodes, episodeNumber)
	if err != nil {
		return nil, "", err
	}
	episodePage, err := ds.HttpWrapper.Get(episodeLink)
	if err != nil {
		return nil, "", err
	}
	defer episodePage.Body.Close()
	linkWithMirror, err := ds.ScrapService.GetLinkWithMirror(episodePage.Body)
	if err != nil {
		return nil, "", err
	}
	pageWithMirror, err := ds.HttpWrapper.Get(linkWithMirror)
	if err != nil {
		return nil, "", err
	}

	downloadUrl, err := ds.ScrapService.GetMegauploadEpisodeLink(pageWithMirror.Body)
	if err != nil {
		return nil, "", err
	}
	log.Println("url is " + downloadUrl)
	return ds.downloadFromM4upload(downloadUrl)

}

func (ds DowloaderService) getLinkForEpisodeNumber(espisodeLinks []string, episodeNumber string) (string, error) {
	for _, episodeLink := range espisodeLinks {
		if ds.ScrapService.GetEpisodeNumber(episodeLink) == episodeNumber {
			return episodeLink, nil
		}
	}
	return "", errors.New("episode number download link not found")
}

func (ds DowloaderService) getEpisodesAvaliable(serieLink string) ([]string, error) {

	episodesPage, err := ds.HttpWrapper.Get(serieLink)
	if err != nil {
		return nil, err
	}
	defer episodesPage.Body.Close()

	episodes, err := ds.ScrapService.GetEpisodesList(episodesPage.Body)
	if err != nil {
		return nil, err
	}
	sort.Strings(episodes)
	return episodes, nil
}

func (ds DowloaderService) GetSortedEpisodesAvaliable(serieLink string) ([]string, error) {

	episdesLinks, err := ds.getEpisodesAvaliable(serieLink)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)

	for _, episodeLink := range episdesLinks {

		result = append(result, ds.ScrapService.GetEpisodeNumber(episodeLink))
	}
	sort.Strings(result)
	return result, nil
}

func (ds DowloaderService) downloadFromM4upload(downloadUrl string) (io.Reader, string, error) {
	code, err := ds.ScrapService.GetMegauploadCode(downloadUrl)
	if err != nil {
		return nil, "", err
	}
	payload := fmt.Sprintf("op=download2&id=%s&method_free=+", code)
	postResult, err := ds.HttpWrapper.Request(downloadUrl, "POST", strings.NewReader(payload))
	if err != nil {
		return nil, "", err
	}
	defer postResult.Body.Close()
	downloadLink := postResult.Header.Get("location")

	headers := make(map[string]string)
	headers["Referer"] = "https://www.mp4upload.com/"

	log.Println("downloading from" + downloadLink)

	episodeResp, err := ds.HttpWrapper.RequestWithHeaders(downloadLink, "GET", nil, headers)
	if err != nil {
		return nil, "", err
	}

	return episodeResp.Body, "mp4", nil
}
