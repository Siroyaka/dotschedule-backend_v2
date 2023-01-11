package rssmaster

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	rssmaster "github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/rss/master"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type UpdateRepository struct {
	sqlHandler                   abstruct.SqlHandler
	updateRssMasterQueryTemplate string
}

func NewUpdateRepository(handler abstruct.SqlHandler, updateRssMasterQueryTemplate string) rssmaster.UpdateRepository {
	return UpdateRepository{
		sqlHandler:                   handler,
		updateRssMasterQueryTemplate: updateRssMasterQueryTemplate,
	}
}

func (r UpdateRepository) UpdateTime(target string, date utility.WrappedTime) utility.IError {
	sqmt, err := r.sqlHandler.Prepare(r.updateRssMasterQueryTemplate)
	if err != nil {
		return utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, r.updateRssMasterQueryTemplate)
	}
	defer sqmt.Close()

	_, err = sqmt.Exec(date.ToUTCFormatString(), target)

	if err != nil {
		return utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}
	return nil
}
