package fullschedule

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type InsertRepository interface {
	Insert(domain.FullScheduleData, utility.WrappedTime) (int64, utility.IError)
}
