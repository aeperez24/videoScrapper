package editor

import (
	"aeperez24/videoScrapper/application"
	"aeperez24/videoScrapper/service"
)

func GetProviders() []string {
	providerList := []string{"animeShow", "cuevana"}
	return providerList

}

func GetConfigurations(configPath string) []service.SerieConfiguration {
	return application.LoadConfigurationWithPath("./").SerieConfigurations
}
