package viewschedule

import (
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetRepository interface {
	GetScheduleData(utility.WrappedTime, utility.WrappedTime, int, ScheduleDataAdapter) ([]domain.ScheduleData, utility.IError)
	GetMonthData(fromDate utility.WrappedTime, toDate utility.WrappedTime, displayScheduleStatus int, dataAdapter MonthDataAdapter) ([]domain.MonthData, utility.IError)
}

type ScheduleDataAdapter func(
	id, platform, url, streamer_name, streamer_id, title, description string,
	status int,
	publish_datetime string,
	duration int,
	thumbnail, icon, participants_data string,
) (domain.ScheduleData, utility.IError)

type MonthDataAdapter func(string, string, string) domain.MonthData
