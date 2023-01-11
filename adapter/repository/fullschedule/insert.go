package fullschedule

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/fullschedule"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type InsertRepository struct {
	sqlHandler                  abstruct.SqlHandler
	insertScheduleQueryTemplate string
}

func NewInsertRepository(
	sqlHandler abstruct.SqlHandler,
	insertScheduleQueryTemplate string,
) fullschedule.InsertRepository {
	return InsertRepository{
		sqlHandler:                  sqlHandler,
		insertScheduleQueryTemplate: insertScheduleQueryTemplate,
	}
}

func (repos InsertRepository) Insert(data domain.FullScheduleData, insertTime utility.WrappedTime) (int64, utility.IError) {
	sqmt, err := repos.sqlHandler.Prepare(repos.insertScheduleQueryTemplate)
	if err != nil {
		return 0, utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, repos.insertScheduleQueryTemplate)
	}
	defer sqmt.Close()

	result, err := sqmt.Exec(data.StreamingID, data.PlatformType, data.Url, data.StreamerName, data.StreamerID, data.Title, data.Description, data.Status, data.PublishDatetime, data.Duration, data.ThumbnailLink, insertTime.ToUTCFormatString(), insertTime.ToUTCFormatString(), data.IsViewing, data.IsCompleteData)
	if err != nil {
		return 0, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return count, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	return count, nil
}
