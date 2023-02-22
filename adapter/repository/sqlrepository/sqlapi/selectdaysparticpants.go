package sqlapi

import (
	"fmt"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain/apidomain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference/apireference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type SelectDaysParticipantsRepository struct {
	selectWrapper       sqlwrapper.SelectWrapper[apireference.DayParticipantsResponse]
	viewStatus          int
	localTimeDifference int
}

func NewSelectDaysParticipantsRepository(
	sqlHandler abstruct.SqlHandler,
	query string,
	viewStatus int,
	localTimeDifference int,
) SelectDaysParticipantsRepository {
	return SelectDaysParticipantsRepository{
		selectWrapper:       sqlwrapper.NewSelectWrapper[apireference.DayParticipantsResponse](sqlHandler, query),
		viewStatus:          viewStatus,
		localTimeDifference: localTimeDifference,
	}
}

func (repos SelectDaysParticipantsRepository) scan(s sqlwrapper.IScan) (apireference.DayParticipantsResponse, utilerror.IError) {
	var date string
	var particpants string
	if err := s.Scan(
		&date,
		&particpants,
	); err != nil {
		return apireference.DayParticipantsResponse{}, utilerror.New(err.Error(), "")
	}

	participantsData, ierr := utility.JsonUnmarshal[[]apidomain.StreamerData](particpants)
	if ierr != nil {
		participantsData = []apidomain.StreamerData{}
		fmt.Println("partse error")
		fmt.Println(ierr.Error())
	}

	return apireference.DayParticipantsResponse{
		Date:         date,
		Participants: participantsData,
	}, nil
}

func (repos SelectDaysParticipantsRepository) Execute(date apireference.FromToDate) ([]apireference.DayParticipantsResponse, utilerror.IError) {
	dateFrom := date.From.ToLocalFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())
	dateTo := date.To.ToLocalFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())
	result, err := repos.selectWrapper.SelectPrepare(repos.scan, repos.localTimeDifference, dateFrom, dateTo, repos.viewStatus, repos.localTimeDifference)
	if err != nil {
		return result, err.WrapError()
	}
	return result, nil
}
