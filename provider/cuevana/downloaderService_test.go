package cuevana

import (
	"aeperez24/animewatcher/service"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCuevanaService = DownloaderService{ScrapService: ScrapperService{}, HttpWrapper: buildHttpWrapperMock(),
	getProxies: mockGetProxies, getClientWithProxies: mockGetClientWithProxy, usedProxies: make(map[string]bool)}

func TestMustReturnSortedEpisodesLinks(t *testing.T) {
	links, _ := getTestCuevanaService().GetSortedEpisodesLinks("serieLink")
	assert.Equal(t, 30, len(links))
	assert.True(t, contains(links, "https://ww3.cuevana3.me/episodio/servant-1x1"))
}

func TestMustReturnEpisodesNames(t *testing.T) {
	episodes, _ := getTestCuevanaService().GetSortedEpisodesAvaliable("serieLink")
	assert.Equal(t, 30, len(episodes))
	assert.True(t, contains(episodes, "servant-1x1"))
}

func TestMustReturnDownloadEpisodeFromLink(t *testing.T) {
	video, format, err := getTestCuevanaService().DownloadEpisodeFromLink("serieLink", "servant-1x1")
	assert.Nil(t, err)
	bytes, _ := ioutil.ReadAll(video)
	assert.Equal(t, "video", string(bytes))
	assert.Equal(t, "mp4", format)

}

func getTestCuevanaService() DownloaderService {
	return DownloaderService{ScrapService: ScrapperService{}, HttpWrapper: buildHttpWrapperMock(),
		getProxies: mockGetProxies, getClientWithProxies: mockGetClientWithProxy, usedProxies: make(map[string]bool)}
}
func mockGetProxies() []string {
	return []string{"proxy"}
}

func mockGetClientWithProxy(in []string) (httpPostClient, string) {
	return buildHttpWrapperMock(), in[0]
}

func buildHttpWrapperMock() service.HttpWrapper {
	httpWrapper := service.HttpWrapperMock{}

	httpWrapper.On("Get", "serieLink").Return(&http.Response{
		Body: open("inputs/seriepage.html"),
	}, nil)

	httpWrapper.On("Get", "https://ww3.cuevana3.me/episodio/servant-1x1").Return(&http.Response{
		Body: open("inputs/episodepage.html"),
	}, nil)

	httpWrapper.On("Get", "https://1fichier.com/?o0oflhfdfby481t7e3bg#Synchronization+Service").Return(&http.Response{
		Body: open("inputs/1fichier.html"),
	}, nil)

	httpWrapper.On("PostForm", "https://1fichier.com/?o0oflhfdfby481t7e3bg",
		url.Values{"adz": []string{"313.776397987367"}, "dl_no_ssl": []string{"on"}, "dlinline": []string{"on"}}).
		Return(&http.Response{
			Body: open("inputs/downloadpage.html"),
		}, nil)

	httpWrapper.On("Get", "http://a-11.1fichier.com/c262664672?inline").
		Return(&http.Response{
			Body: io.NopCloser(strings.NewReader("video")),
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
