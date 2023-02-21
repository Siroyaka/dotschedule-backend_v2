package interactor

import (
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference/apireference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type DaysParticipantsInteractor struct {
	selectRepository abstruct.RepositoryRequest[apireference.FromToDate, []apireference.DayParticipantsResponse]
}

func NewDaysParticipantsInteractor(
	selectRepository abstruct.RepositoryRequest[apireference.FromToDate, []apireference.DayParticipantsResponse],
) DaysParticipantsInteractor {
	return DaysParticipantsInteractor{
		selectRepository: selectRepository,
	}
}

func (intr DaysParticipantsInteractor) GetMonthData(fromDate, toDate wrappedbasics.IWrappedTime) ([]apireference.DayParticipantsResponse, utility.IError) {
	result, err := intr.selectRepository.Execute(apireference.FromToDate{From: fromDate, To: toDate})
	if err != nil {
		return result, err.WrapError()
	}

	return result, nil
}
