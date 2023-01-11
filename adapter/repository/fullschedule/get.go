package fullschedule

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/fullschedule"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetRepository[X any] struct {
	sqlHandler abstruct.SqlHandler
	query      string
}

func NewGetRepository[X any](handler abstruct.SqlHandler, query string) fullschedule.GetRepository[X] {
	return GetRepository[X]{
		sqlHandler: handler,
		query:      query,
	}
}

func (repos GetRepository[X]) Get(f func(utility.IScan) (X, error), data ...any) (result []X, resErr utility.IError) {
	sqmt, err := repos.sqlHandler.Prepare(repos.query)
	if err != nil {
		resErr = utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, repos.query)
		return
	}
	defer sqmt.Close()

	rows, err := sqmt.Query(data...)
	if err != nil {
		resErr = utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
		return
	}
	defer rows.Close()

	for rows.Next() {
		val, err := f(rows)
		if err != nil {
			utility.LogError(utility.NewError(err.Error(), utility.ERR_SQL_DATASCAN))
			continue
		}
		result = append(result, val)
	}
	resErr = nil
	return
}
