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
	//TODO
	ds := service.DowloaderService{
		ScrapService:     service.ScrapperServiceImpl{},
		Tracker:          service.TrackerServiceImpl{},
		GetSender:        service.GetWrapper{},
		FileSystemSaver:  service.FileSystemSaverWrapper{},
		AppConfiguration: appConfig,
	}

	for i, config := range appConfig.AnimeConfigurations {
		go asyncDownload(ds.DownloadLastEpisode, config.AnimeLink, chanArr[i])
	}

	for _, channel := range chanArr {
		err = <-channel
		if err != nil {
			log.Println(err)
		} else {
			log.Println("a download completed")
		}

	}

}

func asyncDownload(fn func(string) (string, error), in string, errorChanel chan error) {
	_, err := fn(in)
	errorChanel <- err
}
