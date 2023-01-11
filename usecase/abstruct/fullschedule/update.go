package fullschedule

import (
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type UpdateAnyColumnRepository interface {
	Update(utility.WrappedTime, ...any) (int64, utility.IError)
}
