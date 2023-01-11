package streamingparticipants

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type DeleteRepository interface {
	Delete(domain.StreamingParticipants) (int64, utility.IError)
}
