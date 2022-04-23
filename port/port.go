package port

import "io"

type GeneralDownloadService interface {
	GetSortedEpisodesAvaliable(serieLink string) ([]string, error)
	DownloadEpisodeFromLink(serieLink string, episodeNumber string) (io.Reader, error)
}
