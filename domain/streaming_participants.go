package domain

import (
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type StreamingParticipants struct {
	streamingID string
	platform    string
	list        utility.HashSet[string]
}

func EmptyStreamingParticipants() StreamingParticipants {
	return NewStreamingParticipants("", "")
}

func NewStreamingParticipants(streamingID, platform string, list ...string) StreamingParticipants {
	p := StreamingParticipants{streamingID: streamingID, platform: platform, list: utility.NewHashSet[string]()}
	for _, v := range list {
		p = p.Add(v)
	}
	return p
}

func (p StreamingParticipants) StreamingID() string {
	return p.streamingID
}

func (p StreamingParticipants) Platform() string {
	return p.platform
}

func (p StreamingParticipants) Add(value string) StreamingParticipants {
	p.list.Set(value)
	return p
}

func (p StreamingParticipants) Has(value string) bool {
	return p.list.Has(value)
}

func (p StreamingParticipants) GetList() []string {
	return p.list.List()
}

func (p StreamingParticipants) IsEmpty() bool {
	if p.streamingID == "" {
		return true
	}
	return !p.list.Any()
}

type PlatformParticipants struct {
	platformParticipantsData StreamingParticipants
	convertList              map[string]string
}

func NewPlatformParticipants(streamingID, platform string, list ...string) PlatformParticipants {
	return PlatformParticipants{
		platformParticipantsData: NewStreamingParticipants(streamingID, platform, list...),
		convertList:              make(map[string]string),
	}
}

func (pp PlatformParticipants) AddConvertData(streamerID, platformID string) PlatformParticipants {
	pp.convertList[platformID] = streamerID
	return pp
}

func (pp PlatformParticipants) Convert() StreamingParticipants {
	streamerIDSet := utility.NewHashSet[string]()
	for _, v := range pp.platformParticipantsData.list.List() {
		streamerIDSet.Set(pp.convertList[v])
	}
	return StreamingParticipants{
		streamingID: pp.platformParticipantsData.streamingID,
		platform:    pp.platformParticipantsData.platform,
		list:        streamerIDSet,
	}
}
