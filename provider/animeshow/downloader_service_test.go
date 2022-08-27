package animeshow

import (
	serviceMock "aeperez24/animewatcher/mock/service"
	"aeperez24/animewatcher/service"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const EPISODE_LINK = "https://www2.animeshow.tv/Fullmetal-Alchemist-Brotherhood-episode-1/"
const SERIE_LINK = "https://www2.animeshow.tv/Fullmetal-Alchemist-Brotherhood/"

const DOWNLOAD_URL = "https://www.mp4upload.com/a6xkfdysqdbu"
const LOCATION_URL = "location"
const VIDEO_DOWNLOAD = "videoDownload"
const MIRROR_LINK = "https://www2.animeshow.tv/Tate-no-Yuusha-no-Nariagari-Season-2-episode-2-mirror-3/"

func TestGetEpisodesAvaliable(t *testing.T) {
	ds := DowloaderService{
		ScrapService:     ScrapperServiceImpl{},
		HttpWrapper:      buildHttpWrapperMock(),
		AppConfiguration: service.AppConfiguration{},
	}
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
	episodesList, _ := ds.GetSortedEpisodesAvaliable(SERIE_LINK)
	assert.Len(t, episodesList, 64, "the size expected is 64")
	assert.Equal(t, "1", episodesList[0], "")
}

func TestDownloadFromM4upload(t *testing.T) {
	ds := DowloaderService{
		ScrapService:     ScrapperServiceImpl{},
		HttpWrapper:      buildHttpWrapperMock(),
		AppConfiguration: service.AppConfiguration{},
	}

	reader, _, _ := ds.downloadFromM4upload(DOWNLOAD_URL)
	result, _ := ioutil.ReadAll(reader)
	assert.Equal(t, VIDEO_DOWNLOAD, string(result))
}

func TestDownloadFromLink(t *testing.T) {
	ds := DowloaderService{
		ScrapService:     ScrapperServiceImpl{},
		HttpWrapper:      buildHttpWrapperMock(),
		AppConfiguration: service.AppConfiguration{},
	}
	reader, _, _ := ds.DownloadEpisodeFromLink(SERIE_LINK, "1")
	result, _ := ioutil.ReadAll(reader)
	assert.Equal(t, VIDEO_DOWNLOAD, string(result))
}

func buildHttpWrapperMock() service.HttpWrapper {
	httpWrapper := serviceMock.HttpWrapper{}
	httpWrapper.On("Get", SERIE_LINK).Return(&http.Response{
		Body: open("inputs/episodesList.html"),
	}, nil)

	httpWrapper.On("Get", EPISODE_LINK).Return(&http.Response{
		Body: open("inputs/episode.html"),
	}, nil)

	httpWrapper.On("Get", MIRROR_LINK).Return(&http.Response{
		Body: open("inputs/episodeMU.html"),
	}, nil)

	headersLocation := http.Header{}
	headersLocation["Location"] = []string{LOCATION_URL}
	httpWrapper.On("Request", DOWNLOAD_URL, "POST", mock.Anything).Return(&http.Response{
		Header: headersLocation,
		Body:   io.NopCloser(strings.NewReader("")),
	}, nil)

	headersRefer := make(map[string]string)
	headersRefer["Referer"] = "https://www.mp4upload.com/"
	httpWrapper.On("RequestWithHeaders", LOCATION_URL, "GET", nil, headersRefer).Return(&http.Response{
		Body: io.NopCloser(strings.NewReader(VIDEO_DOWNLOAD)),
	}, nil)

	return &httpWrapper
}

func open(filepath string) *os.File {
	res, _ := os.Open(filepath)
	return res
}
