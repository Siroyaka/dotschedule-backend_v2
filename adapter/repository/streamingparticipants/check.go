package streamingparticipants

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type CheckRepository struct {
	sqlHandler    abstruct.SqlHandler
	queryTemplate string
}

func NewCheckRepository(sqlHandler abstruct.SqlHandler, queryTemplate string) CheckRepository {
	return CheckRepository{
		sqlHandler:    sqlHandler,
		queryTemplate: queryTemplate,
	}
}

func (repos CheckRepository) Check(streamingId, platform, member_id string, cntRange int) (bool, utility.IError) {
	sqmt, err := repos.sqlHandler.Prepare(repos.queryTemplate)
	if err != nil {
		return false, utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, repos.queryTemplate)
	}
	defer sqmt.Close()

	row, err := sqmt.Query(streamingId, platform)
	if err != nil {
		return false, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	result := false
	for row.Next() {
		var count int
		row.Scan(&count)
		result = count > cntRange
	}
	return result, nil
}
