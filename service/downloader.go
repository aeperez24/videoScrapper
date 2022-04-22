package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
)

type AnimeConfiguration struct {
	AnimeLink string
	AnimeName string
	Provider  string
}
type AppConfiguration struct {
	AnimeConfigurations []AnimeConfiguration
	OutputPath          string
}

type GeneralDownloadService interface {
	GetSortedEpisodesAvaliable(serieLink string) ([]string, error)
	DownloadEpisodeFromLink(serieLink string, episodeNumber string) (io.Reader, error)
}

type DownloaderManager struct {
	FileSystemManager  FileSystemManager
	AppConfiguration   AppConfiguration
	DownloaderServices map[string]GeneralDownloadService
	Tracker            TrackerService
}

func (dm DownloaderManager) getAnimeNameFromLink(link string) string {
	animeConfgs := dm.AppConfiguration.AnimeConfigurations
	for _, config := range animeConfgs {
		if config.AnimeLink == link {
			return config.AnimeName
		}
	}
	return ""
}

func (dm DownloaderManager) getConfigFromLink(link string) (AnimeConfiguration, error) {
	animeConfgs := dm.AppConfiguration.AnimeConfigurations
	for _, config := range animeConfgs {
		if config.AnimeLink == link {
			return AnimeConfiguration{}, nil
		}
	}
	return AnimeConfiguration{}, errors.New(fmt.Sprintf("configuratio not found for %v", link))
}
func (dm DownloaderManager) DownloadLastEpisode(animeLink string) (string, error) {
	animeConfig, err := dm.getConfigFromLink(animeLink)
	if err != nil {
		return "", err
	}
	downloadService, ok := dm.DownloaderServices[animeConfig.Provider]
	if !ok {
		return "", errors.New(fmt.Sprintf("downloader not found for %v", animeConfig.Provider))
	}

	episodesAvaliable, err := downloadService.GetSortedEpisodesAvaliable(animeLink)

	if err != nil {
		return "", err
	}

	lastEpisodeAvaliable := episodesAvaliable[len(episodesAvaliable)-1]
	isDownloaded := dm.Tracker.IsPreviouslyDownloaded(animeLink, lastEpisodeAvaliable)
	if isDownloaded {
		log.Default().Printf("episode %v for link %v is already downloaded", lastEpisodeAvaliable, animeLink)
	}
	episodeReader, err := downloadService.DownloadEpisodeFromLink(animeLink, lastEpisodeAvaliable)
	if err != nil {
		return "", err
	}
	dm.FileSystemManager.Save(dm.AppConfiguration.OutputPath+"/"+animeConfig.AnimeName, lastEpisodeAvaliable, episodeReader)
	dm.Tracker.SaveAlreadyDownloaded(animeLink, lastEpisodeAvaliable)
	return "", nil
}

type DowloaderService struct {
	ScrapService     ScrapperService
	GetSender        GetSender
	AppConfiguration AppConfiguration
}

func (dls DowloaderService) DownloadEpisodeFromLink(serieLink string, episodeNumber string) (io.Reader, error) {
	episodes, err := dls.GetSortedEpisodesAvaliable(serieLink)

	if err != nil {
		return nil, err
	}

	episodeLink, _ := dls.getLinkForEpisodeNumber(episodes, episodeNumber)

	episodePage, _ := dls.GetSender.Get(episodeLink)
	defer episodePage.Body.Close()
	linkWithMirror, _ := dls.ScrapService.GetLinkWithMirror(episodePage.Body)
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

func (dls DowloaderService) GetSortedEpisodesAvaliable(serieLink string) ([]string, error) {

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
