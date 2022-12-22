package application

import (
	"aeperez24/videoScrapper/port"
	"aeperez24/videoScrapper/provider/animeshow"
	"aeperez24/videoScrapper/provider/cuevana"
	"aeperez24/videoScrapper/service"
	"log"
	"os"
)

const defaultOutput = "/output/"

type Application struct {
	downloadServicesMap map[string]port.GeneralDownloadService
	configuration       service.AppConfiguration
	downloadManager     service.DownloaderManager
}

func NewApplication() Application {
	appConfig := loadConfiguration()
	if appConfig.OutputPath == "" {
		appConfig.OutputPath = defaultOutput
	}

	configureLogs(appConfig)
	servicesMap := initializeDownloadServices(appConfig)
	downloaderManager := service.DownloaderManagerImpl{FileSystemManager: service.FileSystemManagerWrapper{},
		AppConfiguration: appConfig, DownloaderServices: servicesMap, Tracker: service.TrackerServiceImpl{FileSystemManager: service.FileSystemManagerWrapper{}},
	}
	return Application{
		servicesMap,
		appConfig,
		downloaderManager,
	}

}

func loadConfiguration() service.AppConfiguration {
	return LoadConfigurationWithPath("./")
}

func LoadConfigurationWithPath(configPath string) service.AppConfiguration {
	appConfig, err := service.LoadConfig(configPath)

	if err != nil {
		log.Fatal(err)
		log.Panic(err)
	}
	return *appConfig
}
func configureLogs(appConfig service.AppConfiguration) {
	if appConfig.LogsPath != "" {
		logFile, err := os.OpenFile(appConfig.LogsPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Panic(err)
		}
		log.SetOutput(logFile)
		defer logFile.Close()
	}
}

func initializeDownloadServices(appConfig service.AppConfiguration) map[string]port.GeneralDownloadService {
	dsAnimeshow := animeshow.DowloaderService{
		ScrapService:     animeshow.ScrapperServiceImpl{},
		HttpWrapper:      service.HttpWrapperImpl{},
		AppConfiguration: appConfig,
	}

	dsCuevana := cuevana.NewDownloaderService(service.HttpWrapperImpl{})
	servicesMap := map[string]port.GeneralDownloadService{}
	servicesMap["animeshow"] = dsAnimeshow
	servicesMap["cuevana"] = dsCuevana
	return servicesMap
}

func (app Application) Run() {
	chanArr := make([]chan []error, len(app.configuration.SerieConfigurations))
	for i := range chanArr {
		chanArr[i] = make(chan []error)
	}
	for i, config := range app.configuration.SerieConfigurations {
		go asyncDownload(app.downloadManager.DownloadAllEpisodes, config.SerieLink, chanArr[i])
	}

	for _, channel := range chanArr {
		errList := <-channel
		printErrorList(errList)
	}

}

func asyncDownload(fn func(string) []error, in string, errorChanel chan []error) {
	defer close(errorChanel)
	err := fn(in)
	errorChanel <- err
}

func printErrorList(errList []error) {
	for _, err := range errList {
		log.Println(err)
	}
}
