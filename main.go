package main

import (
	"aeperez24/animewatcher/service"
	"log"
)

func main() {

	appConfig, err := service.LoadConfig("./")

	if err != nil {
		log.Fatal(err)
		log.Panic(err)
	}

	chanArr := make([]chan error, len(appConfig.AnimeConfigurations))
	for i, _ := range chanArr {
		chanArr[i] = make(chan error)
	}
	ds := service.DowloaderService{
		ScrapService:      service.ScrapperServiceImpl{},
		Tracker:           service.TrackerServiceImpl{FileSystemManager: service.FileSystemManagerWrapper{}},
		GetSender:         service.GetWrapper{},
		FileSystemManager: service.FileSystemManagerWrapper{},
		AppConfiguration:  appConfig,
	}

	for i, config := range appConfig.AnimeConfigurations {
		go asyncDownload(ds.DownloadLastEpisode, config.AnimeLink, chanArr[i])
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
