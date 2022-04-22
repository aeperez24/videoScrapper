package service

import (
	"fmt"
	"log"
	"strings"
)

const TRACKING_FILES_PATH = "tracking_files"

type TrackerService interface {
	IsPreviouslyDownloaded(animeLink string, episodelink string) bool
	SaveAlreadyDownloaded(animeLink string, episodelink string)
}

type TrackerServiceImpl struct {
	FileSystemManager FileSystemManager
}

func (trackerService TrackerServiceImpl) IsPreviouslyDownloaded(animeName string, episodeNumber string) bool {
	byteArr, _ := trackerService.FileSystemManager.Read(TRACKING_FILES_PATH, animeName)
	if byteArr == nil {
		return false
	}
	stringFile := fmt.Sprintf("%s", byteArr)
	episodesNumber := strings.Split(stringFile, " ")
	log.Println(episodesNumber)
	for _, ep := range episodesNumber {
		if ep == episodeNumber {
			return true
		}
	}
	return false
}

func (trackerService TrackerServiceImpl) SaveAlreadyDownloaded(animeName string, episodeNumber string) {
	byteArr, _ := trackerService.FileSystemManager.Read(TRACKING_FILES_PATH, animeName)
	if byteArr == nil {
		byteArr = []byte("")
	}
	stringFile := fmt.Sprintf("%s", byteArr)
	stringFile = stringFile + " " + episodeNumber
	reader := strings.NewReader(stringFile)
	trackerService.FileSystemManager.Save(TRACKING_FILES_PATH, animeName, reader)

}
