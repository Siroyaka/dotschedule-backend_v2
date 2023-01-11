package migration

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetDataRepository interface {
	GetStreamerMaster(string) (map[string]domain.GroupStreamerData, utility.IError)
	GetIDSet([]string, string) (utility.HashSet[string], utility.IError)
}
