package cuevana

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ScrapperServiceImpl struct{}

func (ScrapperServiceImpl) getEpisodesList(data io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return nil, err
	}
	episodesArr := make([]string, 0)
	doc.Find(".TpRwCont .all-episodes  article a").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		episodesArr = append(episodesArr, val)
	})
	fmt.Println(episodesArr)
	return episodesArr, nil
}

func (ScrapperServiceImpl) getEpisodeName(link string) string {
	episodeNumber := (strings.Split(link, "/episodio/"))[1]
	episodeNumber = strings.ReplaceAll(episodeNumber, "/", "")
	return episodeNumber
}

func (ScrapperServiceImpl) get1fichierLink(data io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return "", err
	}
	resultLink := ""
	doc.Find(".mdl-bd .TPTblCn a").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		if strings.Contains(val, "1fichier.com") {
			resultLink = val
		}

	})
	if resultLink == "" {
		err = errors.New("1fichier.com link not found")
	}

	return resultLink, err
}
