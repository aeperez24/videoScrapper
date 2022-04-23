package port

import (
	"aeperez24/animewatcher/service"
	"io"
)

type GeneralDownloadService interface {
	GetSortedEpisodesAvaliable(serieLink string) ([]string, error)
	DownloadEpisodeFromLink(serieLink string, episodeNumber string) (io.Reader, error)
}

type DowloaderService struct {
	ScrapService     ScrapperService
	GetSender        service.GetSender
	AppConfiguration service.AppConfiguration
}
