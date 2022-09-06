package service

import (
	"aeperez24/videoScrapper/port"
	"errors"
	"fmt"
)

type SerieConfiguration struct {
	SerieLink string
	SerieName string
	Provider  string
}
type AppConfiguration struct {
	SerieConfigurations []SerieConfiguration
	OutputPath          string
	LogsPath            string
}

type DownloaderManagerImpl struct {
	FileSystemManager  FileSystemManager
	AppConfiguration   AppConfiguration
	DownloaderServices map[string]port.GeneralDownloadService
	Tracker            TrackerService
}

type DownloaderManager interface {
	DownloadLastEpisode(SerieLink string) []error
	DownloadAllEpisodes(SerieLink string) []error
}

func (dm DownloaderManagerImpl) getSerieNameFromLink(link string) string {
	animeConfgs := dm.AppConfiguration.SerieConfigurations
	for _, config := range animeConfgs {
		if config.SerieLink == link {
			return config.SerieName
		}
	}
	return ""
}

func (dm DownloaderManagerImpl) getConfigFromLink(link string) (SerieConfiguration, error) {
	animeConfgs := dm.AppConfiguration.SerieConfigurations
	for _, config := range animeConfgs {
		if config.SerieLink == link {
			return config, nil
		}
	}
	return SerieConfiguration{}, errors.New(fmt.Sprintf("configuratio not found for %v", link))
}
func (dm DownloaderManagerImpl) DownloadLastEpisode(SerieLink string) []error {
	animeConfig, err := dm.getConfigFromLink(SerieLink)
	if err != nil {
		return []error{err}
	}
	downloadService, ok := dm.DownloaderServices[animeConfig.Provider]
	if !ok {
		return []error{errors.New(fmt.Sprintf("downloader not found for %v", animeConfig.Provider))}
	}
	episodesAvaliable, err := downloadService.GetSortedEpisodesAvaliable(SerieLink)
	if err != nil {
		return []error{err}
	}

	lastEpisodeAvaliable := episodesAvaliable[len(episodesAvaliable)-1]
	_, err = dm.downloadEpisode(SerieLink, lastEpisodeAvaliable, animeConfig, downloadService)
	if err != nil {
		return []error{err}
	}
	return []error{}
}

func (dm DownloaderManagerImpl) downloadEpisode(serieLink string, episodeNumber string,
	serieConfig SerieConfiguration, downloadService port.GeneralDownloadService) (string, error) {

	isDownloaded := dm.Tracker.IsPreviouslyDownloaded(serieConfig.SerieName, episodeNumber)
	if isDownloaded {
		return "", nil
	}
	episodeReader, format, err := downloadService.DownloadEpisodeFromLink(serieLink, episodeNumber)
	if err != nil {
		return "", err
	}
	err = dm.FileSystemManager.Save(dm.AppConfiguration.OutputPath+"/"+serieConfig.SerieName, episodeNumber+"."+format, episodeReader)
	if err != nil {
		return "", err
	}
	dm.Tracker.SaveAlreadyDownloaded(serieConfig.SerieName, episodeNumber)
	return "", nil
}

func (dm DownloaderManagerImpl) DownloadAllEpisodes(SerieLink string) []error {
	animeConfig, err := dm.getConfigFromLink(SerieLink)
	if err != nil {
		return []error{err}
	}
	downloadService, ok := dm.DownloaderServices[animeConfig.Provider]
	if !ok {
		return []error{errors.New(fmt.Sprintf("downloader not found for %v", animeConfig.Provider))}
	}
	episodesAvaliable, err := downloadService.GetSortedEpisodesAvaliable(SerieLink)
	if err != nil {
		return []error{err}
	}
	errorList := make([]error, 0)
	for _, episode := range episodesAvaliable {
		_, err := dm.downloadEpisode(SerieLink, episode, animeConfig, downloadService)
		if err != nil {
			errorList = append(errorList, err)
		}
	}
	return errorList
}
