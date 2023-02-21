package apireference

import "github.com/Siroyaka/dotschedule-backend_v2/domain/apidomain"

type ScheduleResponse struct {
	StreamingData apidomain.StreamingData
	StreamerData  apidomain.StreamerData
	Participants  []apidomain.StreamerData
}
