package editor

import (
	"aeperez24/videoScrapper/application"
	"aeperez24/videoScrapper/service"
	"fmt"
	"os"
)

const EDITOR_CONFIG_PATH = "./"
const EDITOR_CONFIG_NAME = "editor"
const CONFIG_PATH_VAR_NAME = "WCH_PATH"

func GetProviders() []string {
	providerList := []string{"animeshow", "cuevana"}
	return providerList

}

func GetConfigurations(configPath string) []service.SerieConfiguration {
	return application.LoadConfigurationWithPath("./").SerieConfigurations
}

func GetConfigPath() (string, error) {
	path := os.Getenv(CONFIG_PATH_VAR_NAME)
	if path != "" {
		return path, nil
	}
	return "", fmt.Errorf("env var %v not found", CONFIG_PATH_VAR_NAME)
}
