package test

import (
	"aeperez24/animewatcher/service"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEpisodesList(t *testing.T) {
	htmlFile, err := os.Open("inputs/episodesList.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlFile.Close()

	scraperService := service.ScrapperServiceImpl{}

	episodesList, _ := scraperService.GetEpisodesList(htmlFile)
	assert.Len(t, episodesList, 64, "the size expected is 64")
	assert.Equal(t, "https://www2.animeshow.tv/Fullmetal-Alchemist-Brotherhood-episode-64/", episodesList[0], "")
}

func TestGetMegaupLoadEpisodeLink(t *testing.T) {
	htmlFile, err := os.Open("inputs/episodeMu.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlFile.Close()

	scraperService := service.ScrapperServiceImpl{}

	episodeLink, _ := scraperService.GetMegauploadEpisodeLink(htmlFile)
	assert.Equal(t, "https://www.mp4upload.com/embed-a6xkfdysqdbu.html", episodeLink, "error getting downLoadLink")
}
