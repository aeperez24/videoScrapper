package cuevana

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEpisodesList(t *testing.T) {
	htmlFile, err := os.Open("inputs/seriepage.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlFile.Close()

	scraperService := ScrapperServiceImpl{}

	episodesList, _ := scraperService.GetEpisodesList(htmlFile)
	assert.Len(t, episodesList, 30, "the size expected is 30")
	assert.Equal(t, "https://ww3.cuevana3.me/episodio/servant-1x1", episodesList[0], "")
}
