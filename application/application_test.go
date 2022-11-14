package application

import (
	serviceMock "aeperez24/videoScrapper/mock/service"
	"aeperez24/videoScrapper/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	dm := serviceMock.DownloaderManager{}
	dm.Mock.On("DownloadAllEpisodes", "testSerieLink1").Return([]error{})
	dm.Mock.On("DownloadAllEpisodes", "testSerieLink2").Return([]error{})

	app := Application{
		downloadServicesMap: nil,
		configuration: service.AppConfiguration{
			SerieConfigurations: []service.SerieConfiguration{
				{SerieName: "testSerieName1", SerieLink: "testSerieLink1", Provider: "testProvider1"},
				{SerieName: "testSerieName2", SerieLink: "testSerieLink2", Provider: "testProvider2"},
			},
		}, downloadManager: &dm,
	}
	app.Run()
	dm.Mock.AssertCalled(t, "DownloadAllEpisodes", "testSerieLink1")
	dm.Mock.AssertCalled(t, "DownloadAllEpisodes", "testSerieLink2")
	assert.Equal(t, true, true)
}

func TestLoadConfig(t *testing.T) {
	config := LoadConfigurationWithPath("../")
	assert.Equal(t, "/Users/path/downloads", config.OutputPath)
	assert.Equal(t, "cuevana|animeshow", config.SerieConfigurations[0].Provider)
	assert.Equal(t, "serieLink", config.SerieConfigurations[0].SerieLink)
	assert.Equal(t, "serieName", config.SerieConfigurations[0].SerieName)

}

func TestInitializeDownloadServicesl(t *testing.T) {
	expectedServices := []string{"animeshow", "cuevana"}
	config := LoadConfigurationWithPath("../")
	serviceMap := initializeDownloadServices(config)
	for _, service := range expectedServices {
		_, isContained := serviceMap[service]
		assert.True(t, isContained)
	}
	assert.Equal(t, len(expectedServices), len(serviceMap))

}
