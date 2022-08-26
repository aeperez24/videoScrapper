package animeshow

import (
	serviceMock "aeperez24/animewatcher/mock/service"
	"aeperez24/animewatcher/service"
	"log"
	"net/http"
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

const SERIE_LINK = "serieLink"

func TestGetEpisodesAvaliable(t *testing.T) {
	ds := DowloaderService{
		ScrapService:     ScrapperServiceImpl{},
		HttpWrapper:      buildHttpWrapperMock(),
		AppConfiguration: service.AppConfiguration{},
	}
	htmlFile, err := os.Open("inputs/episodesList.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlFile.Close()

	episodesList, _ := ds.getEpisodesAvaliable(SERIE_LINK)
	assert.Len(t, episodesList, 64, "the size expected is 64")
	assert.Equal(t, "https://www2.animeshow.tv/Fullmetal-Alchemist-Brotherhood-episode-1/", episodesList[0], "")

}

func TestGetSortedEpisodesAvaliable(t *testing.T) {
	ds := DowloaderService{
		ScrapService:     ScrapperServiceImpl{},
		HttpWrapper:      buildHttpWrapperMock(),
		AppConfiguration: service.AppConfiguration{},
	}
	htmlFile, err := os.Open("inputs/episodesList.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlFile.Close()

	episodesList, _ := ds.GetSortedEpisodesAvaliable(SERIE_LINK)
	assert.Len(t, episodesList, 64, "the size expected is 64")
	assert.Equal(t, "1", episodesList[0], "")
}

func buildHttpWrapperMock() service.HttpWrapper {
	httpWrapper := serviceMock.HttpWrapper{}
	httpWrapper.On("Get", SERIE_LINK).Return(&http.Response{
		Body: open("inputs/episodesList.html"),
	}, nil)
	return &httpWrapper
}

func open(filepath string) *os.File {
	res, _ := os.Open(filepath)
	return res
}
