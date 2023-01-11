package viewschedule

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/viewschedule"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetRepository struct {
	sqlHandler                abstruct.SqlHandler
	scheduleGetQueryTemplate  string
	monthDataGetQueryTemplate string
	localTimeDifference       string
}

func NewGetRepository(sqlHandler abstruct.SqlHandler, scheduleGetQueryTemplate string, monthDataGetQueryTemplate string, localTimeDifference string) viewschedule.GetRepository {
	return GetRepository{
		sqlHandler:                sqlHandler,
		scheduleGetQueryTemplate:  scheduleGetQueryTemplate,
		monthDataGetQueryTemplate: monthDataGetQueryTemplate,
		localTimeDifference:       localTimeDifference,
	}
}

func (repos GetRepository) GetScheduleData(
	fromDate, toDate utility.WrappedTime,
	displayScheduleStatus int,
	dataAdapter viewschedule.ScheduleDataAdapter,
) ([]domain.ScheduleData, utility.IError) {
	var list []domain.ScheduleData

	stmt, err := repos.sqlHandler.Prepare(repos.scheduleGetQueryTemplate)
	if err != nil {
		return list, utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, repos.scheduleGetQueryTemplate)
	}
	defer stmt.Close()

	rows, err := stmt.Query(fromDate.ToUTCFormatString(), toDate.ToUTCFormatString(), displayScheduleStatus)

	if err != nil {
		return list, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var platform string
		var url string
		var streamer_name string
		var streamer_id string
		var title string
		var description string
		var status int
		var publish_datetime string
		var duration int
		var thumbnail string
		var icon string
		var participants_data string
		err := rows.Scan(&id, &platform, &url, &streamer_name, &streamer_id, &title, &description, &status, &publish_datetime, &duration, &thumbnail, &icon, &participants_data)
		if err != nil {
			utility.LogError(utility.NewError(err.Error(), utility.ERR_SQL_DATASCAN))
			continue
		}
		data, ierr := dataAdapter(id, platform, url, streamer_name, streamer_id, title, description, status, publish_datetime, duration, thumbnail, icon, participants_data)
		if ierr != nil {
			utility.LogError(ierr.WrapError())
			continue
		}
		list = append(list, data)
	}

	return list, nil
}

func (repos GetRepository) GetMonthData(
	fromDate, toDate utility.WrappedTime,
	displayScheduleStatus int,
	dataAdapter viewschedule.MonthDataAdapter,
) ([]domain.MonthData, utility.IError) {
	var list []domain.MonthData

	stmt, err := repos.sqlHandler.Prepare(repos.monthDataGetQueryTemplate)
	if err != nil {
		return list, utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, repos.monthDataGetQueryTemplate)
	}
	defer stmt.Close()

	rows, err := stmt.Query(repos.localTimeDifference, fromDate.ToUTCFormatString(), toDate.ToUTCFormatString(), displayScheduleStatus, repos.localTimeDifference)

	if err != nil {
		return list, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		var memberData string
		var platform_icons string
		err := rows.Scan(&date, &memberData, &platform_icons)
		if err != nil {
			utility.LogError(utility.NewError(err.Error(), utility.ERR_SQL_DATASCAN))
			continue
		}
		data := dataAdapter(date, memberData, platform_icons)
		list = append(list, data)
	}

	return list, nil
}
