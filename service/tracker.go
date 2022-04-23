package service

import (
	"fmt"
	"strings"
)

const TRACKING_FILES_PATH = "tracking_files"

type TrackerService interface {
	IsPreviouslyDownloaded(SerieLink string, episodelink string) bool
	SaveAlreadyDownloaded(SerieLink string, episodelink string)
}

type TrackerServiceImpl struct {
	FileSystemManager FileSystemManager
}

func (trackerService TrackerServiceImpl) IsPreviouslyDownloaded(SerieName string, episodeNumber string) bool {
	byteArr, _ := trackerService.FileSystemManager.Read(TRACKING_FILES_PATH, SerieName)
	if byteArr == nil {
		return false
	}
	stringFile := fmt.Sprintf("%s", byteArr)
	episodesNumber := strings.Split(stringFile, " ")
	for _, ep := range episodesNumber {
		if ep == episodeNumber {
			return true
		}
	}
	return false
}

func (trackerService TrackerServiceImpl) SaveAlreadyDownloaded(SerieName string, episodeNumber string) {
	byteArr, _ := trackerService.FileSystemManager.Read(TRACKING_FILES_PATH, SerieName)
	if byteArr == nil {
		byteArr = []byte("")
	}
	stringFile := fmt.Sprintf("%s", byteArr)
	stringFile = stringFile + " " + episodeNumber
	reader := strings.NewReader(stringFile)
	err := trackerService.FileSystemManager.Save(TRACKING_FILES_PATH, SerieName, reader)
	if err != nil {
		fmt.Println("error saving download track")
		fmt.Println(err)
	}
}
