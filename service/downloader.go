package service

import (
	"errors"
	"fmt"
	"io"
	"log"
)

type SerieConfiguration struct {
	SerieLink string
	SerieName string
	Provider  string
}
type AppConfiguration struct {
	SerieConfigurations []SerieConfiguration
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

func (dm DownloaderManager) getSerieNameFromLink(link string) string {
	animeConfgs := dm.AppConfiguration.SerieConfigurations
	for _, config := range animeConfgs {
		if config.SerieLink == link {
			return config.SerieName
		}
	}
	return ""
}

func (dm DownloaderManager) getConfigFromLink(link string) (SerieConfiguration, error) {
	animeConfgs := dm.AppConfiguration.SerieConfigurations
	for _, config := range animeConfgs {
		if config.SerieLink == link {
			return config, nil
		}
	}
	return SerieConfiguration{}, errors.New(fmt.Sprintf("configuratio not found for %v", link))
}
func (dm DownloaderManager) DownloadLastEpisode(SerieLink string) (string, error) {
	animeConfig, err := dm.getConfigFromLink(SerieLink)
	if err != nil {
		return "", err
	}
	downloadService, ok := dm.DownloaderServices[animeConfig.Provider]
	if !ok {
		return "", errors.New(fmt.Sprintf("downloader not found for %v", animeConfig.Provider))
	}

	episodesAvaliable, err := downloadService.GetSortedEpisodesAvaliable(SerieLink)

	if err != nil {
		return "", err
	}
	lastEpisodeAvaliable := episodesAvaliable[len(episodesAvaliable)-1]
	isDownloaded := dm.Tracker.IsPreviouslyDownloaded(animeConfig.SerieName, lastEpisodeAvaliable)
	if isDownloaded {
		log.Default().Printf("episode %v for link %v is already downloaded", lastEpisodeAvaliable, SerieLink)
		return "", nil
	}
	episodeReader, err := downloadService.DownloadEpisodeFromLink(SerieLink, lastEpisodeAvaliable)
	if err != nil {
		return "", err
	}
	dm.FileSystemManager.Save(dm.AppConfiguration.OutputPath+"/"+animeConfig.SerieName, lastEpisodeAvaliable, episodeReader)
	dm.Tracker.SaveAlreadyDownloaded(SerieLink, lastEpisodeAvaliable)
	return "", nil
}
