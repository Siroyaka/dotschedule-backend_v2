package sqlwrapper

import (
	"fmt"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

type SelectWrapper[X any] struct {
	sqlHandler abstruct.SqlHandler
	query      string
}

type IScan interface {
	Scan(...interface{}) error
}

type Scanable[X any] func(IScan) (X, utilerror.IError)

func NewSelectWrapper[X any](sqlHandler abstruct.SqlHandler, query string) SelectWrapper[X] {
	return SelectWrapper[X]{
		sqlHandler: sqlHandler,
		query:      query,
	}
}

func (repos *SelectWrapper[X]) SetQuery(query string) {
	repos.query = query
}

func (repos SelectWrapper[X]) Select(scanable Scanable[X]) ([]X, utilerror.IError) {
	if repos.query == "" {
		return []X{}, utilerror.New("query is empty.", "")
	}

	rows, err := repos.sqlHandler.Query(repos.query)
	if err != nil {
		return []X{}, utilerror.New(err.Error(), utilerror.ERR_SQL_QUERY, repos.query)
	}
	var result []X
	rowsErrCount := 0
	for rows.Next() {
		res, err := scanable(rows)

		if err != nil {
			logger.Error(err)
			rowsErrCount++
			continue
		}
		result = append(result, res)
	}
	if rowsErrCount != 0 {
		return result, utilerror.New(fmt.Sprintf("scan error: %d", rowsErrCount), utilerror.ERR_SQL_DATASCAN)
	}

	return result, nil
}

func (repos SelectWrapper[X]) SelectPrepare(scanable Scanable[X], data ...interface{}) ([]X, utilerror.IError) {
	if repos.query == "" {
		return []X{}, utilerror.New("query is empty.", "")
	}

	sqmt, err := repos.sqlHandler.Prepare(repos.query)
	if err != nil {
		return []X{}, utilerror.New(err.Error(), utilerror.ERR_SQL_PREPARE, repos.query)
	}

	rows, err := sqmt.Query(data...)
	if err != nil {
		return []X{}, utilerror.New(err.Error(), utilerror.ERR_SQL_QUERY, repos.query)
	}

	var result []X
	rowsErrCount := 0
	for rows.Next() {
		res, err := scanable(rows)
		if err != nil {
			logger.Error(err)
			rowsErrCount++
			continue
		}
		result = append(result, res)
	}
	if rowsErrCount != 0 {
		return result, utilerror.New(fmt.Sprintf("scan error: %d", rowsErrCount), utilerror.ERR_SQL_DATASCAN)
	}

	return result, nil
}
