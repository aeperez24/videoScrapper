package main

import (
	"aeperez24/animewatcher/port"
	"aeperez24/animewatcher/service"
	"aeperez24/animewatcher/vendors/animeshow"
	"aeperez24/animewatcher/vendors/cuevana"
	"crypto/tls"
	"log"
	"net/http"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	appConfig, err := service.LoadConfig("./")

	if err != nil {
		log.Fatal(err)
		log.Panic(err)
	}

	chanArr := make([]chan []error, len(appConfig.SerieConfigurations))
	for i, _ := range chanArr {
		chanArr[i] = make(chan []error)
	}
	ds := animeshow.DowloaderService{
		ScrapService:     animeshow.ScrapperServiceImpl{},
		GetSender:        service.GetWrapper{},
		AppConfiguration: appConfig,
	}
	dsCuevana := cuevana.DownloaderService{ScrapService: cuevana.ScrapperService{},
		GetSender: service.GetWrapper{}}

	servicesMap := map[string]port.GeneralDownloadService{}
	servicesMap["animeshow"] = ds
	servicesMap["cuevana"] = dsCuevana
	downloaderManager := service.DownloaderManager{FileSystemManager: service.FileSystemManagerWrapper{},
		AppConfiguration: appConfig, DownloaderServices: servicesMap, Tracker: service.TrackerServiceImpl{FileSystemManager: service.FileSystemManagerWrapper{}},
	}

	for i, config := range appConfig.SerieConfigurations {
		go asyncDownload(downloaderManager.DownloadAllEpisodes, config.SerieLink, chanArr[i])
	}

	for _, channel := range chanArr {
		errList := <-channel
		printErrorList(errList)
	}

}

func asyncDownload(fn func(string) []error, in string, errorChanel chan []error) {
	err := fn(in)
	errorChanel <- err
}

func printErrorList(errList []error) {
	for _, err := range errList {
		log.Println(err)
	}
}
