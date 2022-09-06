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

const (
	serieLink     = "serieLink"
	serieName     = "serieName"
	outputPath    = "output"
	episodeNumber = "1"
	fileData      = "videoData"
	fileFormat    = "format"
	provider      = "provider"
)

func TestGetSerieNameFromLink(t *testing.T) {
	assert.Equal(t, serieName, buildDownloadManager().getSerieNameFromLink(serieLink))
}
func TestReturnEmptyStringWhenSerieNameFromLinkIsNotFoundInConfigs(t *testing.T) {

	assert.Equal(t, "", buildDownloadManager().getSerieNameFromLink("wrong"))
}

func TestGetConfigFromLink(t *testing.T) {
	config, err := buildDownloadManager().getConfigFromLink(serieLink)
	assert.Nil(t, err)
	assert.Equal(t, getSerieConfiguration(), config)
}

func TestDownloadEpisode(t *testing.T) {
	generalDonwloaderServiceMock := getGeneralDownloaderMock()
	trackerMock := getTrackerMock()
	dm := buildDownloadManager()
	dm.Tracker = &trackerMock
	dm.downloadEpisode(serieLink, episodeNumber, getSerieConfiguration(), &generalDonwloaderServiceMock)

}

func TestDownloadLastEpisode(t *testing.T) {
	assert.Empty(t, buildDownloadManager().DownloadLastEpisode(serieLink))
}

func TestDownloadAlltEpisodes(t *testing.T) {
	assert.Empty(t, buildDownloadManager().DownloadAllEpisodes(serieLink))
}

//testTarget
func buildDownloadManager() DownloaderManagerImpl {
	fileSystemMock := getFileSystemManagerMock()
	trackerMock := getTrackerMock()
	generalDownloaderMock := getGeneralDownloaderMock()
	return DownloaderManagerImpl{
		AppConfiguration:   getAppConfiguration(),
		FileSystemManager:  &fileSystemMock,
		DownloaderServices: map[string]port.GeneralDownloadService{provider: &generalDownloaderMock},
		Tracker:            &trackerMock,
	}
}

//mocks building
func getFileSystemManagerMock() serviceMock.FileSystemManager {
	fileSystemMock := serviceMock.FileSystemManager{}
	fileSystemMock.On("Save", outputPath+"/"+serieName, episodeNumber+"."+fileFormat, io.NopCloser(strings.NewReader(fileData))).Return(nil)
	return fileSystemMock
}

func getTrackerMock() serviceMock.TrackerService {
	trackerMock := serviceMock.TrackerService{}
	trackerMock.On("IsPreviouslyDownloaded", serieName, episodeNumber).Return(false)
	trackerMock.On("SaveAlreadyDownloaded", serieName, episodeNumber).Return()
	return trackerMock

}

func getGeneralDownloaderMock() portMock.GeneralDownloadService {
	generalDonwloaderServiceMock := portMock.GeneralDownloadService{}

	generalDonwloaderServiceMock.On("DownloadEpisodeFromLink", serieLink, episodeNumber).Return(io.NopCloser(strings.NewReader(fileData)), "format", nil)
	generalDonwloaderServiceMock.On("GetSortedEpisodesAvaliable", serieLink).Return([]string{episodeNumber}, nil)

	return generalDonwloaderServiceMock
}

func getAppConfiguration() AppConfiguration {
	return AppConfiguration{
		OutputPath:          outputPath,
		SerieConfigurations: []SerieConfiguration{getSerieConfiguration()},
	}
}

func getSerieConfiguration() SerieConfiguration {
	return SerieConfiguration{SerieLink: serieLink, SerieName: serieName, Provider: provider}
}
