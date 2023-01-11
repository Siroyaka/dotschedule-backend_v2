package rssschedule

import (
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type UpdateRepository interface {
	Update([]string, string, int) utility.IError
}
