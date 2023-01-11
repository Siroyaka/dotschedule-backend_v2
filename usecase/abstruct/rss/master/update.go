package rssmaster

import (
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type UpdateRepository interface {
	UpdateTime(string, utility.WrappedTime) utility.IError
}
