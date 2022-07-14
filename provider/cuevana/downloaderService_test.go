package cuevana

import (
	"aeperez24/animewatcher/service"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustReturnSortedEpisodesLinks(t *testing.T) {
	httpWrapper := buildHttpWrapperMock()
	ds := DownloaderService{ScrapService: ScrapperService{}, HttpWrapper: httpWrapper}
	links, _ := ds.GetSortedEpisodesLinks("serieLink")
	assert.Equal(t, 30, len(links))
	assert.True(t, contains(links, "https://ww3.cuevana3.me/episodio/servant-3x9"))
}

func buildHttpWrapperMock() service.HttpWrapper {
	httpWrapper := service.HttpWrapperMock{}

	httpWrapper.On("Get", "serieLink").Return(&http.Response{
		Body: open("inputs/seriepage.html"),
	}, nil)

	return httpWrapper
}

func open(filepath string) *os.File {
	res, _ := os.Open(filepath)
	return res
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
