package service

type TrackerService interface {
	IsPreviouslyDownloaded(animeLink string, episodelink string) bool
	SaveAlreadyDownloaded(animeLink string, episodelink string)
}

type TrackerServiceImpl struct {
}

//TODO
func (TrackerServiceImpl) IsPreviouslyDownloaded(animeLink string, episodelink string) bool {
	return false
}

//TODO
func (TrackerServiceImpl) SaveAlreadyDownloaded(animeLink string, episodelink string) {

}
