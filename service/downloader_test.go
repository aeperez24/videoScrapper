package service

import (
	"io"
	"strings"
	"testing"

	portMock "aeperez24/videoScrapper/mock/port"
	serviceMock "aeperez24/videoScrapper/mock/service"
	"aeperez24/videoScrapper/port"

	"github.com/stretchr/testify/assert"
)

const serie_link = "serieLink"
const serie_name = "serieName"
const output_path = "output"
const episode_number = "1"
const file_data = "videoData"
const file_format = "format"
const provider = "provider"

func TestGetSerieNameFromLink(t *testing.T) {
	assert.Equal(t, serie_name, buildDownloadManager().getSerieNameFromLink(serie_link))
}
func TestReturnEmptyStringWhenSerieNameFromLinkIsNotFoundInConfigs(t *testing.T) {

	assert.Equal(t, "", buildDownloadManager().getSerieNameFromLink("wrong"))
}

func TestGetConfigFromLink(t *testing.T) {
	config, err := buildDownloadManager().getConfigFromLink(serie_link)
	assert.Nil(t, err)
	assert.Equal(t, getSerieConfiguration(), config)
}

func TestDownloadEpisode(t *testing.T) {
	generalDonwloaderServiceMock := getGeneralDownloaderMock()
	trackerMock := getTrackerMock()
	dm := buildDownloadManager()
	dm.Tracker = &trackerMock
	dm.downloadEpisode(serie_link, episode_number, getSerieConfiguration(), &generalDonwloaderServiceMock)

}

func TestDownloadLastEpisode(t *testing.T) {
	assert.Empty(t, buildDownloadManager().DownloadLastEpisode(serie_link))
}

func TestDownloadAlltEpisodes(t *testing.T) {
	assert.Empty(t, buildDownloadManager().DownloadAllEpisodes(serie_link))
}

func buildDownloadManager() DownloaderManager {
	fileSystemMock := getFileSystemManagerMock()
	trackerMock := getTrackerMock()
	generalDownloaderMock := getGeneralDownloaderMock()
	return DownloaderManager{
		AppConfiguration:   getAppConfiguration(),
		FileSystemManager:  &fileSystemMock,
		DownloaderServices: map[string]port.GeneralDownloadService{provider: &generalDownloaderMock},
		Tracker:            &trackerMock,
	}
}

func getAppConfiguration() AppConfiguration {
	return AppConfiguration{
		OutputPath:          output_path,
		SerieConfigurations: []SerieConfiguration{getSerieConfiguration()},
	}
}

func getSerieConfiguration() SerieConfiguration {
	return SerieConfiguration{SerieLink: serie_link, SerieName: serie_name, Provider: provider}
}

func getTrackerMock() serviceMock.TrackerService {
	trackerMock := serviceMock.TrackerService{}
	trackerMock.On("IsPreviouslyDownloaded", serie_name, episode_number).Return(false)
	trackerMock.On("SaveAlreadyDownloaded", serie_name, episode_number).Return()
	return trackerMock

}

func getGeneralDownloaderMock() portMock.GeneralDownloadService {
	generalDonwloaderServiceMock := portMock.GeneralDownloadService{}

	generalDonwloaderServiceMock.On("DownloadEpisodeFromLink", serie_link, episode_number).Return(io.NopCloser(strings.NewReader(file_data)), "format", nil)
	generalDonwloaderServiceMock.On("GetSortedEpisodesAvaliable", serie_link).Return([]string{episode_number}, nil)

	return generalDonwloaderServiceMock
}

func getFileSystemManagerMock() serviceMock.FileSystemManager {
	fileSystemMock := serviceMock.FileSystemManager{}
	fileSystemMock.On("Save", output_path+"/"+serie_name, episode_number+"."+file_format, io.NopCloser(strings.NewReader(file_data))).Return(nil)
	return fileSystemMock
}
