package service

import (
	"errors"
	"fmt"
	"io"
	"log"
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
	GetSortedEpisodesAvaliable(serieLink string) []string
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

	episodesAvaliable := downloadService.GetSortedEpisodesAvaliable(animeLink)
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
	episodesPage, err := dls.GetSender.Get(animeLink)
	if err != nil {
		return "", err
	}
	defer episodesPage.Body.Close()

	episodes, err := dls.ScrapService.GetEpisodesList(episodesPage.Body)
	if err != nil {
		return "", err
	}
	if len(episodes) == 0 {
		return "", errors.New("there is not any episode avaliable to download")
	}

	lastEpisodeLink := episodes[0]
	episodeNumber := dls.ScrapService.GetEpisodeNumber(lastEpisodeLink)
	animeName := dls.getAnimeNameFromLink(animeLink)
	isAreadyDownloaded := dls.Tracker.IsPreviouslyDownloaded(animeName, episodeNumber)
	if isAreadyDownloaded {
		return lastEpisodeLink, nil
	}

	lastEpisodePage, _ := dls.GetSender.Get(lastEpisodeLink[0:len(lastEpisodeLink)-1] + "-mirror-4")
	defer lastEpisodePage.Body.Close()
	downloadUrl, err := dls.ScrapService.GetMegauploadEpisodeLink(lastEpisodePage.Body)
	if err != nil {
		return "", err
	}
	log.Println("url is " + downloadUrl)
	code, err := dls.ScrapService.GetMegauploadCode(downloadUrl)
	if err != nil {
		return "", err
	}
	payload := fmt.Sprintf("op=download2&id=%s&method_free=+", code)
	postResult, err := dls.GetSender.Request(downloadUrl, "POST", strings.NewReader(payload))
	if err != nil {
		return "", err
	}
	defer postResult.Body.Close()
	downloadLink := postResult.Header.Get("location")

	headers := make(map[string]string)
	headers["Referer"] = "https://www.mp4upload.com/"
	downloadLink = strings.Replace(downloadLink, "www12", "www14", 1)
	log.Println("downloading from" + downloadLink)

	episodeResp, err := dls.GetSender.RequestWithHeaders(downloadLink, "GET", nil, headers)
	if err != nil {
		return "", err
	}

	defer episodeResp.Body.Close()

	log.Println(episodeResp.Header)

	err = dls.FileSystemManager.Save(dls.AppConfiguration.OutputPath+"/"+animeName, episodeNumber+".mp4", episodeResp.Body)
	if err != nil {
		return "", err
	}
	dls.Tracker.SaveAlreadyDownloaded(animeName, episodeNumber)
	return lastEpisodeLink, nil

}

func (dls DowloaderService) getAnimeNameFromLink(link string) string {
	animeConfgs := dls.AppConfiguration.AnimeConfigurations
	for _, config := range animeConfgs {
		if config.AnimeLink == link {
			return config.AnimeName
		}
	}
	return ""
}
