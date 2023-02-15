package reference

type StreamingIDWithPlatformType struct {
	streamingID  string
	platformType string
}

func NewStreamingIDWithPlatformType(
	streamingID string,
	platformType string,
) StreamingIDWithPlatformType {
	return StreamingIDWithPlatformType{
		streamingID:  streamingID,
		platformType: platformType,
	}
}

func (st StreamingIDWithPlatformType) Extract() (string, string) {
	return st.streamingID, st.platformType
}
