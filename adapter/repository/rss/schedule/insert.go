package rssschedule

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	rssschedule "github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/rss/schedule"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type InsertRepository struct {
	sqlHandler          abstruct.SqlHandler
	insertQueryTemplate string
}

func NewInsertRepository(
	sqlHandler abstruct.SqlHandler,
	insertQueryTemplate string,
) rssschedule.InsertRepository {
	return InsertRepository{
		sqlHandler:          sqlHandler,
		insertQueryTemplate: insertQueryTemplate,
	}
}

func (r InsertRepository) Insert(data domain.SeedSchedule) utility.IError {
	now := wrappedbasics.Now(wrappedbasics.WrappedTimeProps.LocalLocation())

	sqmt, err := r.sqlHandler.Prepare(r.insertQueryTemplate)
	if err != nil {
		return utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, r.insertQueryTemplate)
	}
	defer sqmt.Close()

	nowString := now.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())

	_, err = sqmt.Exec(data.GetID(), data.GetPlatformType(), data.GetStatus(), nowString, nowString, data.GetVisibleStatus(), data.GetCompleteStatus())
	if err != nil {
		return utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	return nil
}
