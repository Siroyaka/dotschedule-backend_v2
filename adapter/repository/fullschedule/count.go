package fullschedule

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/fullschedule"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type CountRepository struct {
	sqlHandler abstruct.SqlHandler
	query      string
}

func NewCountRepository(handler abstruct.SqlHandler, query string) fullschedule.CountRepository {
	return CountRepository{
		sqlHandler: handler,
		query:      query,
	}
}

func (repos CountRepository) Count(data ...any) (int, utility.IError) {
	sqmt, err := repos.sqlHandler.Prepare(repos.query)
	if err != nil {
		return 0, utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, repos.query)
	}
	defer sqmt.Close()

	rows, err := sqmt.Query(data...)
	if err != nil {
		return 0, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}
	defer rows.Close()

	for rows.Next() {
		var count int
		err := rows.Scan(&count)
		if err != nil {
			utility.LogError(utility.NewError(err.Error(), utility.ERR_SQL_DATASCAN))
			continue
		}
		return count, nil
	}
	return 0, nil
}
