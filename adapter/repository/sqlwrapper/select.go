package sqlwrapper

import (
	"fmt"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type SelectRepository[X any] struct {
	sqlHandler abstruct.SqlHandler
	query      string
}

func NewSelectRepository[X any](sqlHandler abstruct.SqlHandler, query string) sqlwrapper.SelectRepository[X] {
	return SelectRepository[X]{
		sqlHandler: sqlHandler,
		query:      query,
	}
}

func (repos SelectRepository[X]) Select(scanable sqlwrapper.Scanable[X]) ([]X, utility.IError) {
	rows, err := repos.sqlHandler.Query(repos.query)
	if err != nil {
		return []X{}, utility.NewError(err.Error(), utility.ERR_SQL_QUERY, repos.query)
	}
	var result []X
	rowsErrCount := 0
	for rows.Next() {
		res, err := scanable(rows)

		if err != nil {
			utility.LogError(err)
			rowsErrCount++
			continue
		}
		result = append(result, res)
	}
	if rowsErrCount != 0 {
		return result, utility.NewError(fmt.Sprintf("scan error: %d", rowsErrCount), utility.ERR_SQL_DATASCAN)
	}

	return result, nil
}

func (repos SelectRepository[X]) SelectPrepare(scanable sqlwrapper.Scanable[X], data []interface{}) ([]X, utility.IError) {
	sqmt, err := repos.sqlHandler.Prepare(repos.query)
	if err != nil {
		return []X{}, utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, repos.query)
	}

	rows, err := sqmt.Query(data...)
	if err != nil {
		return []X{}, utility.NewError(err.Error(), utility.ERR_SQL_QUERY, repos.query)
	}

	var result []X
	rowsErrCount := 0
	for rows.Next() {
		res, err := scanable(rows)
		if err != nil {
			utility.LogError(err)
			rowsErrCount++
			continue
		}
		result = append(result, res)
	}
	if rowsErrCount != 0 {
		return result, utility.NewError(fmt.Sprintf("scan error: %d", rowsErrCount), utility.ERR_SQL_DATASCAN)
	}

	return result, nil
}
