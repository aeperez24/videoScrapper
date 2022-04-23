package animeshow

import (
	"aeperez24/animewatcher/service"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
)

type DowloaderService struct {
	ScrapService     ScrapperService
	GetSender        service.GetSender
	AppConfiguration service.AppConfiguration
}

func (dls DowloaderService) DownloadEpisodeFromLink(serieLink string, episodeNumber string) (io.Reader, error) {
	episodes, err := dls.getEpisodesAvaliable(serieLink)

	if err != nil {
		return nil, err
	}

	episodeLink, err := dls.getLinkForEpisodeNumber(episodes, episodeNumber)
	if err != nil {
		return nil, err
	}
	episodePage, _ := dls.GetSender.Get(episodeLink)
	defer episodePage.Body.Close()
	linkWithMirror, err := dls.ScrapService.GetLinkWithMirror(episodePage.Body)
	if err != nil {
		return nil, err
	}
	pageWithMirror, _ := dls.GetSender.Get(linkWithMirror)

	downloadUrl, err := dls.ScrapService.GetMegauploadEpisodeLink(pageWithMirror.Body)
	if err != nil {
		return nil, err
	}
	log.Println("url is " + downloadUrl)
	return dls.downloadFromM4upload(downloadUrl)

}

func (ds DowloaderService) getLinkForEpisodeNumber(espisodeLinks []string, episodeNumber string) (string, error) {
	for _, episodeLink := range espisodeLinks {
		if ds.ScrapService.GetEpisodeNumber(episodeLink) == episodeNumber {
			return episodeLink, nil
		}
	}
	return "", errors.New("episode number download link not found")
}

func (dls DowloaderService) getEpisodesAvaliable(serieLink string) ([]string, error) {

	episodesPage, err := dls.GetSender.Get(serieLink)
	if err != nil {
		return nil, err
	}
	defer episodesPage.Body.Close()

	episodes, err := dls.ScrapService.GetEpisodesList(episodesPage.Body)
	if err != nil {
		return nil, err
	}
	sort.Strings(episodes)
	return episodes, nil
}

func (dls DowloaderService) GetSortedEpisodesAvaliable(serieLink string) ([]string, error) {

	episdesLinks, err := dls.getEpisodesAvaliable(serieLink)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)

	for _, episodeLink := range episdesLinks {

		result = append(result, dls.ScrapService.GetEpisodeNumber(episodeLink))
	}
	sort.Strings(result)
	return result, nil
}

func (dls DowloaderService) downloadFromM4upload(downloadUrl string) (io.Reader, error) {
	code, err := dls.ScrapService.GetMegauploadCode(downloadUrl)
	if err != nil {
		return nil, err
	}
	payload := fmt.Sprintf("op=download2&id=%s&method_free=+", code)
	postResult, err := dls.GetSender.Request(downloadUrl, "POST", strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer postResult.Body.Close()
	downloadLink := postResult.Header.Get("location")

	headers := make(map[string]string)
	headers["Referer"] = "https://www.mp4upload.com/"
	downloadLink = strings.Replace(downloadLink, "www12", "www14", 1)
	log.Println("downloading from" + downloadLink)

	episodeResp, err := dls.GetSender.RequestWithHeaders(downloadLink, "GET", nil, headers)
	if err != nil {
		return nil, err
	}

	return episodeResp.Body, nil
}
