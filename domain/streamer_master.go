package domain

type StreamerMaster struct {
	StreamerID   string
	StreamerName string
	Birthday     string
	Enrollment   string
}

type StreamerMasterWithPlatformData struct {
	StreamerMaster
	PlatformData map[string]StreamerPlatformMaster
}

func NewStreamerMaster(streamerID string) StreamerMaster {
	return StreamerMaster{
		StreamerID: streamerID,
	}
}

func NewStreamerMasterWithPlatformData(streamerID string) StreamerMasterWithPlatformData {
	return StreamerMasterWithPlatformData{
		StreamerMaster: NewStreamerMaster(streamerID),
		PlatformData:   make(map[string]StreamerPlatformMaster),
	}
}

func (sm *StreamerMasterWithPlatformData) PushPlatformData(platformType string, platformID string) bool {
	if _, res := sm.PlatformData[platformType]; res {
		return true
	}

	sm.PlatformData[platformType] = StreamerPlatformMaster{
		StreamerID:   sm.StreamerID,
		PlatformType: platformType,
	}
	return false
}

type StreamerPlatformMaster struct {
	StreamerID   string
	PlatformType string
	PlatformID   string
	PlatformIcon string
	PlatformLink string
}
