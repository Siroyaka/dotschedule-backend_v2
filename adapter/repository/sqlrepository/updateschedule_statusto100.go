package sqlrepository

import (
	"fmt"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type UpdateScheduleStatusTo100Repository struct {
	updateWrapper sqlwrapper.UpdateWrapper
}

func NewUpdateScheduleStatusTo100Repository(sqlHandler abstruct.SqlHandler, query string) UpdateScheduleStatusTo100Repository {
	return UpdateScheduleStatusTo100Repository{
		updateWrapper: sqlwrapper.NewUpdateWrapper(sqlHandler, query),
	}
}

// 前回のステータスが100ではないものについてはis_complete_dataを無条件で0にする
//
// 加えて、insertされた日から1週間はis_complete_dataを0で保持して、再アップロードされる可能性を考慮する
func (repos UpdateScheduleStatusTo100Repository) isComplete(scheduleData domain.FullScheduleData, now wrappedbasics.WrappedTime) bool {
	if scheduleData.Status != "100" {
		return false
	}

	scheduleInsertAt, err := wrappedbasics.NewWrappedTimeFromUTC(scheduleData.InsertAt, wrappedbasics.WrappedTimeProps.DateTimeFormat())
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	completeDateLine := now.Add(0, 0, -7, 0, 0, 0)

	return scheduleInsertAt.Before(completeDateLine)
}

func (repos UpdateScheduleStatusTo100Repository) Execute(scheduleData domain.FullScheduleData) (reference.DBUpdateResponse, utility.IError) {
	now := wrappedbasics.Now()

	var isCompleteValue int
	if repos.isComplete(scheduleData, now) {
		isCompleteValue = 1
		utility.LogDebug("100 status schedule is complete.")
	} else {
		isCompleteValue = 0
	}

	count, id, err := repos.updateWrapper.UpdatePrepare(
		now.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat()),
		100,
		0,
		isCompleteValue,
		scheduleData.StreamingID,
		scheduleData.PlatformType,
	)

	if err != nil {
		return reference.DBUpdateResponse{}, err.WrapError("status update to 100 failed")
	}
	if count == 0 {
		return reference.DBUpdateResponse{Count: 0, Id: 0}, utility.NewError("status update to 100 failed", "")
	}

	return reference.DBUpdateResponse{Count: count, Id: id}, nil
}
