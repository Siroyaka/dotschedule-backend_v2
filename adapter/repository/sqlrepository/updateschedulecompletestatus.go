package sqlrepository

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type UpdateScheduleCompleteStatusRepository struct {
	updateWrapper sqlwrapper.UpdateWrapper
}

func NewUpdateScheduleCompleteStatusRepository(sqlHandler abstruct.SqlHandler, query string) UpdateScheduleCompleteStatusRepository {
	return UpdateScheduleCompleteStatusRepository{
		updateWrapper: sqlwrapper.NewUpdateWrapper(sqlHandler, query),
	}
}

func (repos UpdateScheduleCompleteStatusRepository) Execute(scheduleData domain.FullScheduleData) (reference.DBUpdateResponse, utility.IError) {
	updateAt := wrappedbasics.Now().ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())
	count, id, err := repos.updateWrapper.UpdatePrepare(updateAt, scheduleData.IsCompleteData, scheduleData.StreamingID, scheduleData.PlatformType)

	if err != nil {
		return reference.DBUpdateResponse{Count: count, Id: id}, err.WrapError()
	}

	return reference.DBUpdateResponse{Count: count, Id: id}, nil
}
