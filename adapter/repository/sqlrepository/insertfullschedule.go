package sqlrepository

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type InsertFullScheduleRepository struct {
	updateWrapper sqlwrapper.UpdateWrapper
}

func NewInsertFullScheduleRepository(
	sqlHandler abstruct.SqlHandler,
	insertScheduleQueryTemplate string,
) InsertFullScheduleRepository {
	return InsertFullScheduleRepository{
		updateWrapper: sqlwrapper.NewUpdateWrapper(sqlHandler, insertScheduleQueryTemplate),
	}
}

func (repos InsertFullScheduleRepository) Execute(data domain.FullScheduleData) (reference.DBUpdateResponse, utility.IError) {
	insertAt := wrappedbasics.Now().ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())

	count, id, err := repos.updateWrapper.UpdatePrepare(
		data.StreamingID,
		data.PlatformType,
		data.Url,
		data.StreamerName,
		data.StreamerID,
		data.Title,
		data.Description,
		data.Status,
		data.PublishDatetime,
		data.Duration,
		data.ThumbnailLink,
		insertAt,
		insertAt,
		data.IsViewing,
		data.IsCompleteData,
	)
	if err != nil {
		return reference.DBUpdateResponse{}, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	res := reference.DBUpdateResponse{
		Count: count,
		Id:    id,
	}
	return res, err
}
