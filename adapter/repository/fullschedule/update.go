package fullschedule

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type UpdateRepository struct {
	sqlHandler abstruct.SqlHandler
	query      string
}

func NewUpdateRepository(handler abstruct.SqlHandler, query string) UpdateRepository {
	return UpdateRepository{
		sqlHandler: handler,
		query:      query,
	}
}

func (repos UpdateRepository) Update(schedule domain.FullScheduleData, updateAt utility.WrappedTime) (int64, utility.IError) {

	sqmt, err := repos.sqlHandler.Prepare(repos.query)
	if err != nil {
		return 0, utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, repos.query)
	}
	defer sqmt.Close()

	result, err := sqmt.Exec(schedule.Url, schedule.StreamerName, schedule.StreamerID, schedule.Title, schedule.Description, schedule.Status, schedule.PublishDatetime, schedule.Duration, schedule.ThumbnailLink, updateAt.ToUTCFormatString(), schedule.IsViewing, schedule.IsCompleteData, schedule.StreamingID, schedule.PlatformType)
	if err != nil {
		return 0, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}
	affectedRowCount, err := result.RowsAffected()
	if err != nil {
		return affectedRowCount, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	return affectedRowCount, nil
}
