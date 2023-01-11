package streamermaster

import "github.com/Siroyaka/dotschedule-backend_v2/utility"

type GetPlatformIdRepository interface {
	GetPlatformIdToStreamerId(platform string) (map[string]string, utility.IError)
}
