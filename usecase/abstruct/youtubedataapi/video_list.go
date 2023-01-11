package youtubedataapi

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type VideoListRepository interface {
	IdSearch([]string, []string) ([]domain.YoutubeVideoData, utility.IError)
}
