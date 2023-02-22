package sqlrepository

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/dbmodels"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

// platformidからstreameridを割り出す情報の取得
type SelectPlatformIDWithStreamerIDRepository struct {
	selectWrapper sqlwrapper.SelectWrapper[dbmodels.KeyValue[string, string]]
}

func NewSelectPlatformIDWithStreamerIDRepository(sqlHandler abstruct.SqlHandler, query string) SelectPlatformIDWithStreamerIDRepository {
	return SelectPlatformIDWithStreamerIDRepository{
		selectWrapper: sqlwrapper.NewSelectWrapper[dbmodels.KeyValue[string, string]](sqlHandler, query),
	}
}

func (repos SelectPlatformIDWithStreamerIDRepository) scan(s sqlwrapper.IScan) (dbmodels.KeyValue[string, string], utilerror.IError) {
	var streamer_id, platform_id string
	if err := s.Scan(&streamer_id, &platform_id); err != nil {
		return dbmodels.EmptyKeyValue[string, string](), utilerror.New(err.Error(), utilerror.ERR_SQL_DATASCAN)
	}

	return dbmodels.NewKeyValue(platform_id, streamer_id), nil
}

func (repos SelectPlatformIDWithStreamerIDRepository) Execute(platform string) (map[string]string, utilerror.IError) {
	res := make(map[string]string)

	result, err := repos.selectWrapper.SelectPrepare(repos.scan, platform)

	if err != nil {
		return res, err.WrapError()
	}

	return dbmodels.KeyValueToMap(result), nil
}
