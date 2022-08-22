package animeshow

import (
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

	scraperService := ScrapperServiceImpl{}

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

	scraperService := ScrapperServiceImpl{}

	episodeLink, _ := scraperService.GetMegauploadEpisodeLink(htmlFile)
	assert.Equal(t, "https://www.mp4upload.com/a6xkfdysqdbu", episodeLink, "error getting downLoadLink")
}

func TestGetLinkWithMirror(t *testing.T) {
	htmlFile, err := os.Open("inputs/episode.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlFile.Close()

	scraperService := ScrapperServiceImpl{}
	episodeLink, _ := scraperService.GetLinkWithMirror(htmlFile)
	assert.Equal(t, "https://www2.animeshow.tv/Tate-no-Yuusha-no-Nariagari-Season-2-episode-2-mirror-3/", episodeLink, "error getting mirror")
}
