package fullschedule

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/fullschedule"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type UpdateAnyColumnRepository struct {
	sqlHandler abstruct.SqlHandler
	query      string
}

func NewUpdateAnyColumnRepository(handler abstruct.SqlHandler, query string) fullschedule.UpdateAnyColumnRepository {
	return UpdateAnyColumnRepository{
		sqlHandler: handler,
		query:      query,
	}
}

func (repos UpdateAnyColumnRepository) Update(updateAt utility.WrappedTime, values ...any) (int64, utility.IError) {
	if len(values) == 0 {
		return 0, nil
	}

	sqmt, err := repos.sqlHandler.Prepare(repos.query)
	if err != nil {
		return 0, utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, repos.query)
	}
	defer sqmt.Close()
	result, err := sqmt.Exec(utility.ToInterfaceSlice(updateAt.ToUTCFormatString(), values)...)
	if err != nil {
		return 0, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}
	affectedRowCount, err := result.RowsAffected()
	if err != nil {
		return affectedRowCount, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	return affectedRowCount, nil
}
