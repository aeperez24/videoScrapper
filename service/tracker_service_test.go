package service

import (
	"strings"
	"testing"

	serviceMock "aeperez24/videoScrapper/mock/service"

	"github.com/stretchr/testify/assert"
)

func TestShouldRetrunTrueWhenIsPreviouslyDownloaded(t *testing.T) {

	fsMock := &serviceMock.FileSystemManager{}
	fsMock.On("Read", "tracking_files", "SerieName").Return(
		[]byte("3 1 2"), nil)
	trackerService := TrackerServiceImpl{FileSystemManager: fsMock}
	assert.True(t, trackerService.IsPreviouslyDownloaded("SerieName", "3"))

}

func TestShouldRetrunFalseWhenIsNotPreviouslyDownloaded(t *testing.T) {

	fsMock := &serviceMock.FileSystemManager{}
	fsMock.On("Read", trackingFilesPath, "SerieName").Return(
		[]byte("1 2"), nil)
	trackerService := TrackerServiceImpl{FileSystemManager: fsMock}
	assert.False(t, trackerService.IsPreviouslyDownloaded("SerieName", "3"))

}

func TestShouldSaveEpisode3OnTrackingFile(t *testing.T) {

	fsMock := &serviceMock.FileSystemManager{}
	fsMock.On("Read", trackingFilesPath, "SerieName").Return(
		[]byte("1 2"), nil)
	downloadedList := strings.NewReader("1 2 3")
	fsMock.On("Save", trackingFilesPath, "SerieName", downloadedList).Return(nil)
	trackerService := TrackerServiceImpl{FileSystemManager: fsMock}
	trackerService.SaveAlreadyDownloaded("SerieName", "3")
	fsMock.AssertCalled(t, "Save", trackingFilesPath, "SerieName", downloadedList)
}
