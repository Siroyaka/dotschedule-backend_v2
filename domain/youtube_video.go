package domain

type YoutubeVideoData struct {
	Id                   string
	Snippet              YoutubeSnippet
	ContentDetails       YoutubeContentDetails
	LiveStreamingDetails YoutubeLiveStreamingDetails
	empty                bool
}

func NewYoutubeVideoData(
	id string,
	snippet YoutubeSnippet,
	contentDetails YoutubeContentDetails,
	liveStreamingDetails YoutubeLiveStreamingDetails,
) YoutubeVideoData {
	return YoutubeVideoData{
		Id:                   id,
		Snippet:              snippet,
		ContentDetails:       contentDetails,
		LiveStreamingDetails: liveStreamingDetails,
		empty:                false,
	}
}

func NewEmptyYoutubeVideoData() YoutubeVideoData {
	return YoutubeVideoData{
		empty: true,
	}
}

func (d YoutubeVideoData) IsEmpty() bool {
	return d.empty
}

type YoutubeSnippet struct {
	Title                string
	ChannelID            string
	ChannelName          string
	Description          string
	LiveBroadcastContent string
	PublishAt            string
	Thumbnail            string
	empty                bool
}

func NewYoutubeSnippet(
	title, channelId, channelName, description string,
	liveBroadcastContent, publishAt, thumbnail string,
) YoutubeSnippet {
	return YoutubeSnippet{
		Title:                title,
		ChannelID:            channelId,
		ChannelName:          channelName,
		Description:          description,
		LiveBroadcastContent: liveBroadcastContent,
		PublishAt:            publishAt,
		Thumbnail:            thumbnail,
		empty:                false,
	}
}

func NewEmptyYoutubeSnippet() YoutubeSnippet {
	return YoutubeSnippet{
		empty: true,
	}
}

func (s YoutubeSnippet) IsEmpty() bool {
	return s.empty
}

type YoutubeContentDetails struct {
	Duration string
	empty    bool
}

func NewYoutubeContentDetails(duration string) YoutubeContentDetails {
	return YoutubeContentDetails{
		Duration: duration,
		empty:    false,
	}
}

func NewEmptyYoutubeContentDetails() YoutubeContentDetails {
	return YoutubeContentDetails{
		empty: true,
	}
}

func (c YoutubeContentDetails) IsEmpty() bool {
	return c.empty
}

type YoutubeLiveStreamingDetails struct {
	ActualStartTime, ActualEndTime       string
	ScheduledStartTime, ScheduledEndTime string
	empty                                bool
}

func NewYoutubeLiveStreamingDetails(
	actualStartTime, actualEndTime string,
	scheduledStartTime, scheduledEndTime string,
) YoutubeLiveStreamingDetails {
	return YoutubeLiveStreamingDetails{
		ActualStartTime:    actualStartTime,
		ActualEndTime:      actualEndTime,
		ScheduledStartTime: scheduledStartTime,
		ScheduledEndTime:   scheduledEndTime,
		empty:              false,
	}
}

func NewEmptyYoutubeLiveStreamingDetails() YoutubeLiveStreamingDetails {
	return YoutubeLiveStreamingDetails{
		empty: true,
	}
}

func (d YoutubeLiveStreamingDetails) IsEmpty() bool {
	return d.empty
}
