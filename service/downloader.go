package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type AnimeConfiguration struct {
	AnimeLink string
	AnimeName string
}
type AppConfiguration struct {
	AnimeConfigurations []AnimeConfiguration
	OutputPath          string
}

type DowloaderService struct {
	ScrapService     ScrapperService
	Tracker          TrackerService
	GetSender        GetSender
	DownloadUrl      string
	FileSystemSaver  FileSystemSaver
	AppConfiguration AppConfiguration
}

func (dls DowloaderService) DownloadLastEpisode(animeLink string) (string, error) {
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
	isAreadyDownloaded := dls.Tracker.IsPreviouslyDownloaded(animeLink, lastEpisodeLink)
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
	animeName := dls.getAnimeNameFromLink(animeLink)
	episodeNumber := dls.ScrapService.GetEpisodeNumber(lastEpisodeLink)
	err = dls.FileSystemSaver.Save(dls.AppConfiguration.OutputPath+"/"+animeName, episodeNumber+".mp4", episodeResp.Body)
	if err != nil {
		return "", err
	}
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
