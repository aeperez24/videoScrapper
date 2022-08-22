package service

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldRetrunTrueWhenIsPreviouslyDownloaded(t *testing.T) {

	fsMock := &FileSystemManagerMock{}
	fsMock.On("Read", "tracking_files", "SerieName").Return(
		[]byte("3 1 2"), nil)
	trackerService := TrackerServiceImpl{FileSystemManager: fsMock}
	assert.True(t, trackerService.IsPreviouslyDownloaded("SerieName", "3"))

}

func TestShouldRetrunFalseWhenIsNotPreviouslyDownloaded(t *testing.T) {

	fsMock := &FileSystemManagerMock{}
	fsMock.On("Read", TRACKING_FILES_PATH, "SerieName").Return(
		[]byte("1 2"), nil)
	trackerService := TrackerServiceImpl{FileSystemManager: fsMock}
	assert.False(t, trackerService.IsPreviouslyDownloaded("SerieName", "3"))

}

func TestShouldSaveEpisode3OnTrackingFile(t *testing.T) {

	fsMock := &FileSystemManagerMock{}
	fsMock.On("Read", TRACKING_FILES_PATH, "SerieName").Return(
		[]byte("1 2"), nil)
	downloadedList := strings.NewReader("1 2 3")
	fsMock.On("Save", TRACKING_FILES_PATH, "SerieName", downloadedList).Return(nil)
	trackerService := TrackerServiceImpl{FileSystemManager: fsMock}
	trackerService.SaveAlreadyDownloaded("SerieName", "3")
	fsMock.AssertCalled(t, "Save", TRACKING_FILES_PATH, "SerieName", downloadedList)
}

type FileSystemManagerMock struct {
	mock.Mock
}

func (wrapper *FileSystemManagerMock) Save(filepath string, fileName string, reader io.Reader) error {
	args := wrapper.Called(filepath, fileName, reader)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (wrapper *FileSystemManagerMock) Read(filepath string, fileName string) ([]byte, error) {
	args := wrapper.Called(filepath, fileName)
	var err error
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return args.Get(0).([]byte), err
}
