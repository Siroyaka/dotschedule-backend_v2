package reference

type StreamingIDListWithPlatformID struct {
	streamingIDList []string
	platformID      string
}

func NewStreamingIDListWithPlatformID(list []string, platform string) StreamingIDListWithPlatformID {
	return StreamingIDListWithPlatformID{
		streamingIDList: list,
		platformID:      platform,
	}
}

func (slp StreamingIDListWithPlatformID) Extract() ([]string, string) {
	return slp.streamingIDList, slp.platformID
}
