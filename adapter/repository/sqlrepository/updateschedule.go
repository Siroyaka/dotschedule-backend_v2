package sqlrepository

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

type UpdateScheduleRepository struct {
	updateWrapper sqlwrapper.UpdateWrapper
}

func NewUpdateScheduleRepository(sqlHandler abstruct.SqlHandler, query string) UpdateScheduleRepository {
	return UpdateScheduleRepository{
		updateWrapper: sqlwrapper.NewUpdateWrapper(sqlHandler, query),
	}
}

func (repos UpdateScheduleRepository) Execute(scheduleData domain.FullScheduleData) (reference.DBUpdateResponse, utilerror.IError) {
	count, id, err := repos.updateWrapper.UpdatePrepare(
		scheduleData.Url,
		scheduleData.StreamerName,
		scheduleData.StreamerID,
		scheduleData.Title,
		scheduleData.Description,
		scheduleData.Status,
		scheduleData.PublishDatetime,
		scheduleData.Duration,
		scheduleData.ThumbnailLink,
		scheduleData.UpdateAt,
		scheduleData.IsViewing,
		scheduleData.IsCompleteData,
		scheduleData.StreamingID,
		scheduleData.PlatformType,
	)

	if err != nil {
		return reference.DBUpdateResponse{Count: count, Id: id}, err.WrapError()
	}

	return reference.DBUpdateResponse{Count: count, Id: id}, nil
}
