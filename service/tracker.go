package service

import (
	"fmt"
	"log"
	"strings"
)

const trackingFilesPath = "tracking_files"

type TrackerService interface {
	IsPreviouslyDownloaded(SerieLink string, episodelink string) bool
	SaveAlreadyDownloaded(SerieLink string, episodelink string)
}

type TrackerServiceImpl struct {
	FileSystemManager FileSystemManager
}

func (trackerService TrackerServiceImpl) IsPreviouslyDownloaded(SerieName string, episodeNumber string) bool {
	byteArr, _ := trackerService.FileSystemManager.Read(trackingFilesPath, SerieName)
	if byteArr == nil {
		return false
	}
	log.Printf("serie %s episode(%s) already downloaded \n", SerieName, episodeNumber)
	stringFile := fmt.Sprintf("%s", byteArr)
	episodesNumber := strings.Split(stringFile, " ")
	for _, ep := range episodesNumber {
		if strings.TrimSpace(ep) == strings.TrimSpace(episodeNumber) {
			return true
		}
	}
	return false
}

func (trackerService TrackerServiceImpl) SaveAlreadyDownloaded(SerieName string, episodeNumber string) {
	byteArr, _ := trackerService.FileSystemManager.Read(trackingFilesPath, SerieName)
	if byteArr == nil {
		byteArr = []byte("")
	}
	stringFile := fmt.Sprintf("%s", byteArr)

	stringFile = stringFile + " " + strings.TrimSpace(episodeNumber)
	reader := strings.NewReader(stringFile)
	err := trackerService.FileSystemManager.Save(trackingFilesPath, SerieName, reader)
	if err != nil {
		fmt.Println("error saving download track")
		fmt.Println(err)
	}
}
