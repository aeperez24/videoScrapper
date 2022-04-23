package main

import (
	"aeperez24/animewatcher/service"
	"aeperez24/animewatcher/vendors/animeshow"
	"log"
)

func main() {

	appConfig, err := service.LoadConfig("./")

	if err != nil {
		log.Fatal(err)
		log.Panic(err)
	}

	chanArr := make([]chan error, len(appConfig.SerieConfigurations))
	for i, _ := range chanArr {
		chanArr[i] = make(chan error)
	}
	ds := animeshow.DowloaderService{
		ScrapService:     animeshow.ScrapperServiceImpl{},
		GetSender:        service.GetWrapper{},
		AppConfiguration: appConfig,
	}

	servicesMap := map[string]service.GeneralDownloadService{}
	servicesMap["animeshow"] = ds
	downloaderManager := service.DownloaderManager{FileSystemManager: service.FileSystemManagerWrapper{},
		AppConfiguration: appConfig, DownloaderServices: servicesMap, Tracker: service.TrackerServiceImpl{FileSystemManager: service.FileSystemManagerWrapper{}},
	}

	for i, config := range appConfig.SerieConfigurations {
		go asyncDownload(downloaderManager.DownloadLastEpisode, config.SerieLink, chanArr[i])
	}

	for _, channel := range chanArr {
		err = <-channel
		if err != nil {
			log.Println(err)
		} else {
			log.Println("download completed")
		}

	}

}

func asyncDownload(fn func(string) (string, error), in string, errorChanel chan error) {
	_, err := fn(in)
	errorChanel <- err
}
