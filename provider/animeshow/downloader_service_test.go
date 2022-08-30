package animeshow

import (
	serviceMock "aeperez24/videoScrapper/mock/service"
	"aeperez24/videoScrapper/service"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	episodeLink   = "https://www2.animeshow.tv/Fullmetal-Alchemist-Brotherhood-episode-1/"
	serieLink     = "https://www2.animeshow.tv/Fullmetal-Alchemist-Brotherhood/"
	downloadUrl   = "https://www.mp4upload.com/a6xkfdysqdbu"
	locationUrl   = "location"
	videoDownload = "videoDownload"
	mirrorLink    = "https://www2.animeshow.tv/Tate-no-Yuusha-no-Nariagari-Season-2-episode-2-mirror-3/"
)

func TestGetEpisodesAvaliable(t *testing.T) {
	ds := DowloaderService{
		ScrapService:     ScrapperServiceImpl{},
		HttpWrapper:      buildHttpWrapperMock(),
		AppConfiguration: service.AppConfiguration{},
	}
	episodesList, _ := ds.getEpisodesAvaliable(serieLink)
	assert.Len(t, episodesList, 64, "the size expected is 64")
	assert.Equal(t, "https://www2.animeshow.tv/Fullmetal-Alchemist-Brotherhood-episode-1/", episodesList[0], "")

}

func TestGetSortedEpisodesAvaliable(t *testing.T) {
	ds := DowloaderService{
		ScrapService:     ScrapperServiceImpl{},
		HttpWrapper:      buildHttpWrapperMock(),
		AppConfiguration: service.AppConfiguration{},
	}
	episodesList, _ := ds.GetSortedEpisodesAvaliable(serieLink)
	assert.Len(t, episodesList, 64, "the size expected is 64")
	assert.Equal(t, "1", episodesList[0], "")
}

func TestDownloadFromM4upload(t *testing.T) {
	ds := DowloaderService{
		ScrapService:     ScrapperServiceImpl{},
		HttpWrapper:      buildHttpWrapperMock(),
		AppConfiguration: service.AppConfiguration{},
	}

	reader, _, _ := ds.downloadFromM4upload(downloadUrl)
	result, _ := ioutil.ReadAll(reader)
	assert.Equal(t, videoDownload, string(result))
}

func TestDownloadFromLink(t *testing.T) {
	ds := DowloaderService{
		ScrapService:     ScrapperServiceImpl{},
		HttpWrapper:      buildHttpWrapperMock(),
		AppConfiguration: service.AppConfiguration{},
	}
	reader, _, _ := ds.DownloadEpisodeFromLink(serieLink, "1")
	result, _ := ioutil.ReadAll(reader)
	assert.Equal(t, videoDownload, string(result))
}

func buildHttpWrapperMock() service.HttpWrapper {
	httpWrapper := serviceMock.HttpWrapper{}
	httpWrapper.On("Get", serieLink).Return(&http.Response{
		Body: open("inputs/episodesList.html"),
	}, nil)

	httpWrapper.On("Get", episodeLink).Return(&http.Response{
		Body: open("inputs/episode.html"),
	}, nil)

	httpWrapper.On("Get", mirrorLink).Return(&http.Response{
		Body: open("inputs/episodeMU.html"),
	}, nil)

	headersLocation := http.Header{}
	headersLocation["Location"] = []string{locationUrl}
	httpWrapper.On("Request", downloadUrl, "POST", mock.Anything).Return(&http.Response{
		Header: headersLocation,
		Body:   io.NopCloser(strings.NewReader("")),
	}, nil)

	headersRefer := make(map[string]string)
	headersRefer["Referer"] = "https://www.mp4upload.com/"
	httpWrapper.On("RequestWithHeaders", locationUrl, "GET", nil, headersRefer).Return(&http.Response{
		Body: io.NopCloser(strings.NewReader(videoDownload)),
	}, nil)

	return &httpWrapper
}

func open(filepath string) *os.File {
	res, _ := os.Open(filepath)
	return res
}
