package service

import (
	"aeperez24/animewatcher/port"
	"errors"
	"fmt"
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

type DownloaderManager struct {
	FileSystemManager  FileSystemManager
	AppConfiguration   AppConfiguration
	DownloaderServices map[string]port.GeneralDownloadService
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
func (dm DownloaderManager) DownloadLastEpisode(SerieLink string) []error {
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

func (dm DownloaderManager) downloadEpisode(serieLink string, episodeNumber string,
	serieConfig SerieConfiguration, downloadService port.GeneralDownloadService) (string, error) {

	isDownloaded := dm.Tracker.IsPreviouslyDownloaded(serieConfig.SerieName, episodeNumber)
	if isDownloaded {
		log.Default().Printf("episode %v for link %v is already downloaded", episodeNumber, serieLink)
		return "", nil
	}
	episodeReader, format, err := downloadService.DownloadEpisodeFromLink(serieLink, episodeNumber)
	if err != nil {
		return "", err
	}
	dm.FileSystemManager.Save(dm.AppConfiguration.OutputPath+"/"+serieConfig.SerieName, episodeNumber+"."+format, episodeReader)
	log.Println("saving already downloaded")
	log.Printf("%s : %s", serieLink, episodeNumber)
	dm.Tracker.SaveAlreadyDownloaded(serieConfig.SerieName, episodeNumber)
	return "", nil
}

func (dm DownloaderManager) DownloadAllEpisodes(SerieLink string) []error {
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
		log.Println(episode)
		_, err := dm.downloadEpisode(SerieLink, episode, animeConfig, downloadService)
		if err != nil {
			errorList = append(errorList, err)
		}
	}
	return errorList
}
