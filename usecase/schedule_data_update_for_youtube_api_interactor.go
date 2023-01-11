package usecase

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/fullschedule"
)

type ScheduleDataUpdateForYoutubeInteractor struct {
	getScheduleRepos fullschedule.GetRepository[domain.FullScheduleData]
}
