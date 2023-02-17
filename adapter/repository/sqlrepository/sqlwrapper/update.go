package sqlwrapper

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type UpdateWrapper struct {
	sqlHandler abstruct.SqlHandler
	query      string
}

func NewUpdateWrapper(sqlHandler abstruct.SqlHandler, query string) UpdateWrapper {
	return UpdateWrapper{
		sqlHandler: sqlHandler,
		query:      query,
	}
}

func (repos *UpdateWrapper) SetQuery(query string) {
	repos.query = query
}

func (repos UpdateWrapper) Update() (int64, int64, utility.IError) {
	if repos.query == "" {
		return 0, 0, utility.NewError("query is empty", "")
	}

	result, err := repos.sqlHandler.Exec(repos.query)
	if err != nil {
		return 0, 0, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, 0, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return count, 0, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	return count, id, nil
}

func (repos UpdateWrapper) UpdatePrepare(values ...interface{}) (int64, int64, utility.IError) {
	if repos.query == "" {
		return 0, 0, utility.NewError("query is empty", "")
	}

	sqmt, err := repos.sqlHandler.Prepare(repos.query)
	if err != nil {
		return 0, 0, utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, repos.query)
	}

	result, err := sqmt.Exec(values...)
	if err != nil {
		return 0, 0, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, 0, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return count, 0, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	return count, id, nil
}
