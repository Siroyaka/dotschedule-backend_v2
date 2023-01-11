package youtubedataapi

import (
	"context"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct/youtubedataapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type API struct {
	service *youtube.Service
}

func NewYoutubeDataAPI(developerKey string) youtubedataapi.API {
	ctx := context.Background()

	service, err := youtube.NewService(ctx, option.WithAPIKey(developerKey))
	if err != nil {
		panic(err)
	}
	return API{
		service: service,
	}
}

func (api API) VideosList(part ...string) youtubedataapi.VideosListCall {
	return VideosListCall{*api.service.Videos.List(part)}
}

type VideosListCall struct {
	youtube.VideosListCall
}

func (vlc VideosListCall) Id(id ...string) youtubedataapi.VideosListCall {
	return VideosListCall{*vlc.VideosListCall.Id(id...)}
}

func (vlc VideosListCall) Do() (youtubedataapi.VideosListResponse, error) {
	res, err := vlc.VideosListCall.Do()
	return &VideosListResponse{content: *res, index: -1}, err
}

type VideosListResponse struct {
	content youtube.VideoListResponse
	index   int
}

func (vlr *VideosListResponse) Next() bool {
	vlr.index++
	return len(vlr.content.Items) > vlr.index
}

func (vlr VideosListResponse) Item() youtubedataapi.Video {
	if vlr.content.Items[vlr.index] == nil {
		return nil
	}
	return Video(*vlr.content.Items[vlr.index])
}

func (vlr VideosListResponse) MarshalJson() ([]byte, error) {
	return vlr.content.MarshalJSON()
}

type Video youtube.Video

func (v Video) GetId() string {
	return v.Id
}

func (v Video) GetSnippet() youtubedataapi.Snippet {
	if v.Snippet == nil {
		return nil
	}
	return Snippet(*v.Snippet)
}

func (v Video) GetLiveStreamingDetails() youtubedataapi.LiveStreamingDetails {
	if v.LiveStreamingDetails == nil {
		return nil
	}
	return LiveStreamingDetails(*v.LiveStreamingDetails)
}

func (v Video) GetContentDetails() youtubedataapi.ContentDetails {
	if v.ContentDetails == nil {
		return nil
	}
	return ContentDetails(*v.ContentDetails)
}

type LiveStreamingDetails youtube.VideoLiveStreamingDetails

func (lsd LiveStreamingDetails) GetScheduledStartTime() string {
	return lsd.ScheduledStartTime
}

func (lsd LiveStreamingDetails) GetScheduledEndTime() string {
	return lsd.ScheduledEndTime
}

func (lsd LiveStreamingDetails) GetActualStartTime() string {
	return lsd.ActualStartTime
}

func (lsd LiveStreamingDetails) GetActualEndTime() string {
	return lsd.ActualEndTime
}

type ContentDetails youtube.VideoContentDetails

func (cd ContentDetails) GetDuration() string {
	return cd.Duration
}

type Snippet youtube.VideoSnippet

func (s Snippet) GetTitle() string {
	return s.Title
}

func (s Snippet) GetThumbnail() string {
	return s.Thumbnails.High.Url
}

func (s Snippet) GetChannelId() string {
	return s.ChannelId
}

func (s Snippet) GetChannelTitle() string {
	return s.ChannelTitle
}

func (s Snippet) GetLiveBroadcastContent() string {
	return s.LiveBroadcastContent
}

func (s Snippet) GetPublishAt() string {
	return s.PublishedAt
}

func (s Snippet) GetDescription() string {
	return s.Description
}
