package sqlcontains

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

type ContainsScheduleRepository struct {
	selectWrapper sqlwrapper.SelectWrapper[int]
}

// スケジュール情報の存在確認をする
func NewContainsScheduleRepository(handler abstruct.SqlHandler, query string) ContainsScheduleRepository {
	return ContainsScheduleRepository{
		selectWrapper: sqlwrapper.NewSelectWrapper[int](handler, query),
	}
}

func (repos ContainsScheduleRepository) scan(s sqlwrapper.IScan) (int, utilerror.IError) {
	var count int
	err := s.Scan(&count)
	if err != nil {
		return 0, utilerror.New(err.Error(), utilerror.ERR_SQL_DATASCAN)
	}
	return count, nil
}

func (repos ContainsScheduleRepository) Execute(data reference.StreamingIDWithPlatformType) (bool, utilerror.IError) {
	id, platform := data.Extract()

	counts, err := repos.selectWrapper.SelectPrepare(repos.scan, id, platform)

	if err != nil {
		return false, err.WrapError()
	}

	if len(counts) == 0 {
		return false, utilerror.New("query count is 0", utilerror.ERR_SQL_DATASCAN)
	}

	return counts[0] != 0, nil
}
