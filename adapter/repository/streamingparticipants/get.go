package streamingparticipants

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetRepository struct {
	sqlHandler    abstruct.SqlHandler
	queryTemplate string
}

func NewGetRepository(sqlHandler abstruct.SqlHandler, queryTemplate string) GetRepository {
	return GetRepository{
		sqlHandler:    sqlHandler,
		queryTemplate: queryTemplate,
	}
}

func (r GetRepository) Get(streamingId, platform string) (domain.StreamingParticipants, utility.IError) {
	sqmt, err := r.sqlHandler.Prepare(r.queryTemplate)
	if err != nil {
		return domain.EmptyStreamingParticipants(), utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, r.queryTemplate)
	}
	defer sqmt.Close()

	row, err := sqmt.Query(streamingId, platform)
	if err != nil {
		return domain.EmptyStreamingParticipants(), utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	result := domain.NewStreamingParticipants(streamingId, platform)
	for row.Next() {
		var member_id string
		row.Scan(&member_id)
		result = result.Add(member_id)
	}
	return result, nil
}
