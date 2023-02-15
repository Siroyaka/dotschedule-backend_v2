package participants

import "github.com/Siroyaka/dotschedule-backend_v2/utility"

type SingleInsertData struct {
	streamingId  string
	platformType string
	streamerId   string
	updateAt     utility.WrappedTime
}

func NewSingleInsertData(
	streamingId string,
	platformType string,
	streamerId string,
	updateAt utility.WrappedTime,
) SingleInsertData {
	return SingleInsertData{
		streamingId:  streamingId,
		platformType: platformType,
		streamerId:   streamerId,
		updateAt:     updateAt,
	}
}

func (ns SingleInsertData) Extract() (string, string, string, utility.WrappedTime) {
	return ns.streamingId, ns.platformType, ns.streamerId, ns.updateAt
}
