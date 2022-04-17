package main

import (
	"aeperez24/animewatcher/service"
	"log"
)

func main() {
	//url := "https://www2.animeshow.tv/Fullmetal-Alchemist-Brotherhood"
	url := "https://www2.animeshow.tv/Tate-no-Yuusha-no-Nariagari-Season-2"
	appConfig := service.AppConfiguration{
		OutputPath:          "/Users/andresperez/Documents/downloads",
		AnimeConfigurations: []service.AnimeConfiguration{{AnimeLink: url, AnimeName: "risingShield"}},
	}
	ds := service.DowloaderService{
		ScrapService:     service.ScrapperServiceImpl{},
		Tracker:          service.TrackerServiceImpl{},
		GetSender:        service.GetWrapper{},
		DownloadUrl:      "https://goload.pro/download?id=%s",
		FileSystemSaver:  service.FileSystemSaverWrapper{},
		AppConfiguration: appConfig,
	}
	_, err := ds.DownloadLastEpisode(url)

	if err != nil {
		log.Fatal(err)
		log.Panic(err)
	}

}
