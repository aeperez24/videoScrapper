package cuevana

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ScrapperService struct{}

func (ScrapperService) getEpisodesList(data io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return nil, err
	}
	episodesArr := make([]string, 0)
	doc.Find(".TpRwCont .all-episodes  article a").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		episodesArr = append(episodesArr, val)
	})
	return episodesArr, nil
}

func (ScrapperService) getEpisodeName(link string) string {
	episodeNumber := (strings.Split(link, "/episodio/"))[1]
	episodeNumber = strings.ReplaceAll(episodeNumber, "/", "")
	return episodeNumber
}

func (ScrapperService) get1fichierLink(data io.Reader) (string, error) {
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

func (service ScrapperService) getParammsForFichierDownload(data io.Reader) (fichierParams, error) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return fichierParams{}, err
	}

	adz, err := service.getAdzFromFichier(doc)
	if err != nil {
		return fichierParams{}, err
	}
	postUrl, err := service.getPostUrlFromFichier(doc)
	if err != nil {
		return fichierParams{}, err
	}

	return fichierParams{adz, postUrl}, nil
}

func (ScrapperService) getAdzFromFichier(doc *goquery.Document) (string, error) {

	selection := doc.Find(`[name="adz"]`).First()
	result, exist := selection.Attr("value")
	if !exist {
		return "", errors.New("adz not found")
	}
	return result, nil
}

func (ScrapperService) getPostUrlFromFichier(doc *goquery.Document) (string, error) {
	selection := doc.Find(`[method="post"]`).First()
	result, exist := selection.Attr("action")
	if !exist {
		return "", errors.New("url not found")
	}
	fmt.Println(exist)
	return result, nil
}

type fichierParams struct {
	adz     string
	postUrl string
}

func (service ScrapperService) getDownloadLink(data io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return "", err
	}
	selection := doc.Find(`[class="ok btn-general btn-orange"]`).First()
	link, exist := selection.Attr("href")
	if !exist {
		return "", errors.New("downloadLink not found")
	}

	return link, nil
}
