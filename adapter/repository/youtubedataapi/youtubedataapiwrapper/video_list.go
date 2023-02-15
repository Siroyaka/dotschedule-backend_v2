package youtubedataapiwrapper

import (
	abstructYoutube "github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct/youtubedataapi"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type VideoListWrapper struct {
	api abstructYoutube.API
}

func NewVideoListWrapper(api abstructYoutube.API) VideoListWrapper {
	return VideoListWrapper{
		api: api,
	}
}

func convertSnippet(snippet abstructYoutube.Snippet) domain.YoutubeSnippet {
	if snippet == nil {
		return domain.NewEmptyYoutubeSnippet()
	}
	return domain.NewYoutubeSnippet(
		snippet.GetTitle(),
		snippet.GetChannelId(),
		snippet.GetChannelTitle(),
		snippet.GetDescription(),
		snippet.GetLiveBroadcastContent(),
		snippet.GetPublishAt(),
		snippet.GetThumbnail(),
	)
}

func convertLiveStreamingDetails(liveStreamingDetails abstructYoutube.LiveStreamingDetails) domain.YoutubeLiveStreamingDetails {
	if liveStreamingDetails == nil {
		return domain.NewEmptyYoutubeLiveStreamingDetails()
	}
	return domain.NewYoutubeLiveStreamingDetails(
		liveStreamingDetails.GetActualStartTime(),
		liveStreamingDetails.GetActualEndTime(),
		liveStreamingDetails.GetScheduledStartTime(),
		liveStreamingDetails.GetScheduledEndTime(),
	)
}

func convertContentDetails(contentDetails abstructYoutube.ContentDetails) domain.YoutubeContentDetails {
	if contentDetails == nil {
		return domain.NewEmptyYoutubeContentDetails()
	}
	return domain.NewYoutubeContentDetails(
		contentDetails.GetDuration(),
	)
}

func (repos VideoListWrapper) IdSearch(part, idList []string) ([]domain.YoutubeVideoData, utility.IError) {
	var resList []domain.YoutubeVideoData
	if len(part) == 0 || len(idList) == 0 {
		return resList, nil
	}

	videoData, err := repos.api.VideosList(part...).Id(idList...).Do()
	if err != nil {
		return resList, utility.NewError(err.Error(), "")
	}

	for {
		if !videoData.Next() {
			break
		}
		item := videoData.Item()

		if item == nil {
			resList = append(resList, domain.NewEmptyYoutubeVideoData())
			continue
		}

		id := item.GetId()

		snippet := convertSnippet(item.GetSnippet())
		liveStreamingDetails := convertLiveStreamingDetails(item.GetLiveStreamingDetails())
		contentDetails := convertContentDetails(item.GetContentDetails())

		video := domain.NewYoutubeVideoData(
			id,
			snippet,
			contentDetails,
			liveStreamingDetails,
		)
		resList = append(resList, video)
	}
	return resList, nil
}
