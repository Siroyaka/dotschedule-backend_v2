package youtubedataapi

type Snippet interface {
	GetTitle() string
	GetDescription() string
	GetThumbnail() string
	GetChannelId() string
	GetChannelTitle() string
	GetLiveBroadcastContent() string
	GetPublishAt() string
}

type ContentDetails interface {
	GetDuration() string
}

type LiveStreamingDetails interface {
	GetScheduledStartTime() string
	GetScheduledEndTime() string
	GetActualStartTime() string
	GetActualEndTime() string
}

type Video interface {
	GetId() string
	GetSnippet() Snippet
	GetLiveStreamingDetails() LiveStreamingDetails
	GetContentDetails() ContentDetails
}

type VideosListResponse interface {
	Next() bool
	Item() Video
	MarshalJson() ([]byte, error)
}

type VideosListCall interface {
	Do() (VideosListResponse, error)
	Id(...string) VideosListCall
}

type API interface {
	VideosList(part ...string) VideosListCall
}
