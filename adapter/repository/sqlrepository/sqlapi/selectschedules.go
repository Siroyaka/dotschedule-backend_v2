package sqlapi

import (
	"fmt"
	"strconv"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain/apidomain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference/apireference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type SelectSchedulesRepository struct {
	selectWrapper sqlwrapper.SelectWrapper[apireference.ScheduleResponse]
	viewStatus    int
}

func NewSelectSchedulesRepository(
	sqlHandler abstruct.SqlHandler,
	query string,
	viewStatus int,
) SelectSchedulesRepository {
	return SelectSchedulesRepository{
		selectWrapper: sqlwrapper.NewSelectWrapper[apireference.ScheduleResponse](sqlHandler, query),
		viewStatus:    viewStatus,
	}
}

func (repos SelectSchedulesRepository) scan(s sqlwrapper.IScan) (apireference.ScheduleResponse, utilerror.IError) {
	var streaming_id string
	var url string
	var platform string
	var status string
	var publish_datetime string
	var duration int
	var thumbnail string
	var title string
	var description string
	var streamer_name string
	var streamer_id string
	var streamer_icon string
	var streamer_link string
	var participants_data string
	if err := s.Scan(
		&streaming_id,
		&url,
		&platform,
		&status,
		&publish_datetime,
		&duration,
		&thumbnail,
		&title,
		&description,
		&streamer_name,
		&streamer_id,
		&streamer_icon,
		&streamer_link,
		&participants_data,
	); err != nil {
		return apireference.ScheduleResponse{}, utilerror.New(err.Error(), "")
	}

	statusNum, err := strconv.Atoi(status)
	if err != nil {
		logger.Error(err)
		return apireference.ScheduleResponse{}, utilerror.New(err.Error(), "")
	}

	startDate, ierr := wrappedbasics.NewWrappedTimeFromUTC(publish_datetime, wrappedbasics.WrappedTimeProps.DateTimeFormat())
	if ierr != nil {
		return apireference.ScheduleResponse{}, ierr.WrapError()
	}

	streamingData := apidomain.StreamingData{
		ID:          streaming_id,
		URL:         url,
		Platform:    platform,
		Status:      statusNum,
		StartDate:   startDate.ToLocalFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat()),
		Duration:    duration,
		Thumbnail:   thumbnail,
		Title:       title,
		Description: description,
	}

	streamerData := apidomain.StreamerData{
		Name:     streamer_name,
		ID:       streamer_id,
		Icon:     streamer_icon,
		Link:     streamer_link,
		Platform: platform,
	}

	participants, ierr := utility.JsonUnmarshal[[]apidomain.StreamerData](participants_data)
	if ierr != nil {
		participants = []apidomain.StreamerData{}
		fmt.Println("partse error")
		fmt.Println(ierr.Error())
	}

	return apireference.ScheduleResponse{
		StreamingData: streamingData,
		StreamerData:  streamerData,
		Participants:  participants,
	}, nil
}

func (repos SelectSchedulesRepository) Execute(date apireference.FromToDate) ([]apireference.ScheduleResponse, utilerror.IError) {
	dateFrom := date.From.ToLocalFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())
	dateTo := date.To.ToLocalFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())
	result, err := repos.selectWrapper.SelectPrepare(repos.scan, dateFrom, dateTo, repos.viewStatus)
	if err != nil {
		return result, err.WrapError()
	}
	return result, nil
}
