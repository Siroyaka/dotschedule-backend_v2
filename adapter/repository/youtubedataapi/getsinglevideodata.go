package youtubedataapi

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct/youtubedataapi"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/youtubedataapi/youtubedataapiwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetSingleVideoDataRepository struct {
	wrapper youtubedataapiwrapper.VideoListWrapper
	part    []string
}

func NewGetSingleVideoDataRepository(api youtubedataapi.API, partList []string) GetSingleVideoDataRepository {
	return GetSingleVideoDataRepository{
		wrapper: youtubedataapiwrapper.NewVideoListWrapper(api),
		part:    partList,
	}
}

func (repos GetSingleVideoDataRepository) Execute(streamingId string) (domain.YoutubeVideoData, utility.IError) {
	videoDataList, err := repos.wrapper.IdSearch(repos.part, []string{streamingId})

	if err != nil {
		return domain.NewEmptyYoutubeVideoData(), err.WrapError()
	}

	if len(videoDataList) == 0 {
		return domain.NewEmptyYoutubeVideoData(), nil
	}

	return videoDataList[0], nil
}
