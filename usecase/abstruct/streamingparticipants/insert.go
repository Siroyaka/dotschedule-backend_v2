package streamingparticipants

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type InsertRepository interface {
	InsertList(streamingId, platformType string, memberIdList ...string) (int64, utility.IError)
	InsertStreamingParticipants(domain.StreamingParticipants) (int64, utility.IError)
}
