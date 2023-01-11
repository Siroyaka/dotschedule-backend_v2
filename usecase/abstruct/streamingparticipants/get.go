package streamingparticipants

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetRepository interface {
	Get(string, string) (domain.StreamingParticipants, utility.IError)
}
