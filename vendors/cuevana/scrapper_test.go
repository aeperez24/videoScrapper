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

	episodesList, _ := scraperService.getEpisodesList(htmlFile)
	assert.Len(t, episodesList, 30, "the size expected is 30")
	assert.Equal(t, "https://ww3.cuevana3.me/episodio/servant-1x1", episodesList[0], "")
}

func TestGetEpisodeNumber(t *testing.T) {
	scraperService := ScrapperServiceImpl{}
	episodeName := scraperService.getEpisodeName("https://ww3.cuevana3.me/episodio/servant-1x1")
	assert.Equal(t, "servant-1x1", episodeName)

}

func TestGet1fichierLink(t *testing.T) {
	//"https://ww3.cuevana3.me/episodio/servant-3x1"
	htmlFile, err := os.Open("inputs/episodepage.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlFile.Close()
	scraperService := ScrapperServiceImpl{}

	ficherLink, _ := scraperService.get1fichierLink(htmlFile)
	assert.Equal(t, "https://1fichier.com/?o0oflhfdfby481t7e3bg#Synchronization+Service", ficherLink)
}
