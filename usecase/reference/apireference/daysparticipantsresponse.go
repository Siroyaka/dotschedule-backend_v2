package apireference

import "github.com/Siroyaka/dotschedule-backend_v2/domain/apidomain"

type DayParticipantsResponse struct {
	Date         string
	Participants []apidomain.StreamerData
}
