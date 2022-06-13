package cuevana

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

type ScrapperServiceImpl struct{}

func (ScrapperServiceImpl) GetEpisodesList(data io.Reader) ([]string, error) {
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
