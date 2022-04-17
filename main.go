package main

import (
	"aeperez24/animewatcher/service"
	"log"
)

func main() {

	appConfig, err := service.LoadConfig("./")
	url := appConfig.AnimeConfigurations[0].AnimeLink
	if err != nil {
		log.Fatal(err)
		log.Panic(err)
	}

	ds := service.DowloaderService{
		ScrapService:     service.ScrapperServiceImpl{},
		Tracker:          service.TrackerServiceImpl{},
		GetSender:        service.GetWrapper{},
		FileSystemSaver:  service.FileSystemSaverWrapper{},
		AppConfiguration: appConfig,
	}
	_, err = ds.DownloadLastEpisode(url)

	if err != nil {
		log.Fatal(err)
		log.Panic(err)
	}

}
