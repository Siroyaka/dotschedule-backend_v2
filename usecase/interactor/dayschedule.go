package interactor

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference/apireference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type DayScheduleInteractor struct {
	getRepos abstruct.RepositoryRequest[apireference.FromToDate, []apireference.ScheduleResponse]
}

func NewDayScheduleInteractor(
	getRepos abstruct.RepositoryRequest[apireference.FromToDate, []apireference.ScheduleResponse],
) DayScheduleInteractor {
	return DayScheduleInteractor{
		getRepos: getRepos,
	}
}

func (intr DayScheduleInteractor) GetScheduleData(fromDate, toDate wrappedbasics.IWrappedTime) ([]apireference.ScheduleResponse, utility.IError) {
	result, err := intr.getRepos.Execute(apireference.FromToDate{From: fromDate, To: toDate})
	if err != nil {
		return result, err.WrapError()
	}
	return result, nil
}

func (intr DayScheduleInteractor) ToJson(data []apireference.ScheduleResponse) (string, utility.IError) {
	d := domain.NewAPIResponseData("ok", len(data), "", data)
	result, err := d.ToJson()
	if err != nil {
		return result, err.WrapError()
	}
	return result, nil
}
