package sqlwrapper

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type UpdateRepository struct {
	sqlHandler abstruct.SqlHandler
	query      string
}

func NewUpdateRepository(sqlHandler abstruct.SqlHandler, query string) sqlwrapper.UpdateRepository {
	return UpdateRepository{
		sqlHandler: sqlHandler,
		query:      query,
	}
}

func (repos UpdateRepository) Update() (int64, int64, utility.IError) {
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

func (repos UpdateRepository) UpdatePrepare(values []interface{}) (int64, int64, utility.IError) {
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
