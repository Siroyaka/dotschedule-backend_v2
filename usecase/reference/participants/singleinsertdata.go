package participants

import (
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type SingleInsertData struct {
	streamingId  string
	platformType string
	streamerId   string
	updateAt     wrappedbasics.IWrappedTime
}

func NewSingleInsertData(
	streamingId string,
	platformType string,
	streamerId string,
	updateAt wrappedbasics.IWrappedTime,
) SingleInsertData {
	return SingleInsertData{
		streamingId:  streamingId,
		platformType: platformType,
		streamerId:   streamerId,
		updateAt:     updateAt,
	}
}

func (ns SingleInsertData) Extract() (string, string, string, wrappedbasics.IWrappedTime) {
	return ns.streamingId, ns.platformType, ns.streamerId, ns.updateAt
}
