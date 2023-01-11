package rssschedule

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	rssschedule "github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/rss/schedule"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type InsertRepository struct {
	sqlHandler          abstruct.SqlHandler
	common              utility.Common
	insertQueryTemplate string
}

func NewInsertRepository(
	sqlHandler abstruct.SqlHandler,
	common utility.Common,
	insertQueryTemplate string,
) rssschedule.InsertRepository {
	return InsertRepository{
		sqlHandler:          sqlHandler,
		common:              common,
		insertQueryTemplate: insertQueryTemplate,
	}
}

func (r InsertRepository) Insert(data domain.SeedSchedule) utility.IError {
	now, ierr := r.common.Now()
	if ierr != nil {
		return ierr.WrapError()
	}
	sqmt, err := r.sqlHandler.Prepare(r.insertQueryTemplate)
	if err != nil {
		return utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, r.insertQueryTemplate)
	}
	defer sqmt.Close()
	_, err = sqmt.Exec(data.GetID(), data.GetPlatformType(), data.GetStatus(), now.ToUTCFormatString(), now.ToUTCFormatString(), data.GetVisibleStatus(), data.GetCompleteStatus())
	if err != nil {
		return utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	return nil
}
