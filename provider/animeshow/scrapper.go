package animeshow

import (
	"errors"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type DownloadLink struct {
	Link    string
	Quality string
}

type ScrapperService interface {
	GetEpisodesList(data io.Reader) ([]string, error)
	GetMegauploadEpisodeLink(data io.Reader) (string, error)
	GetMegauploadCode(string) (string, error)
	GetEpisodeNumber(link string) string
	GetLinkWithMirror(data io.Reader) (string, error)
}

type ScrapperServiceImpl struct{}

func (ScrapperServiceImpl) GetEpisodesList(data io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return nil, err
	}
	episodesArr := make([]string, 0)
	doc.Find(".e_l_r a ").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		episodesArr = append(episodesArr, val)
	})
	return episodesArr, nil
}

func (ScrapperServiceImpl) GetEpisodeNumber(link string) string {
	episodeNumber := (strings.Split(link, "-episode-"))[1]
	episodeNumber = strings.ReplaceAll(episodeNumber, "/", "")
	return episodeNumber
}

func (ScrapperServiceImpl) GetMegauploadEpisodeLink(data io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return "", err
	}
	selection := doc.Find(".embed-responsive-item").First()
	src, exist := selection.Attr("src")
	if !exist {
		return "", errors.New("error src in iframe not found")
	}
	src = strings.Replace(src, "embed-", "", -1)
	return strings.Replace(src, ".html", "", -1), nil
}

func (ScrapperServiceImpl) GetMegauploadCode(uri string) (string, error) {
	splitedUri := strings.Split(uri, "https://www.mp4upload.com/")
	if len(splitedUri) != 2 {
		return "", errors.New("error geting megaupload code ")
	}
	return strings.Replace(splitedUri[1], ".html", "", -1), nil
}

func (ScrapperServiceImpl) GetLinkWithMirror(data io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return "", err
	}

	//log.Println(doc.Find(":contains(MP4UPLOAD) .episode_mirrors_hd :contains(HD)").First().Parent().Siblings().Attr("href"))
	result, _ := doc.Find("a :contains(MP4UPLOAD) .episode_mirrors_hd :contains(HD)").Parent().Parent().Parent().Parent().First().Attr("href")
	if result == "" {
		return result, errors.New("link for m4upload mirror  hd not found")
	}
	return result, nil

}
