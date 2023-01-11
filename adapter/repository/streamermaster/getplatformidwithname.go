package streamermaster

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetPlatformIdWithNameRepository struct {
	sqlHandler abstruct.SqlHandler
	query      string
}

func NewGetPlatformIdWithNameRepository(sqlHandler abstruct.SqlHandler, query string) GetPlatformIdRepository {
	return GetPlatformIdRepository{
		sqlHandler: sqlHandler,
		query:      query,
	}
}

func (repos GetPlatformIdWithNameRepository) GetPlatformIdToStreamerIdWithName(platform string) (map[string]string, utility.IError) {
	res := make(map[string]string)

	stmt, err := repos.sqlHandler.Prepare(repos.query)
	if err != nil {
		return res, utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, repos.query)
	}
	defer stmt.Close()

	rows, err := stmt.Query(platform)
	if err != nil {
		return res, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}
	defer rows.Close()
	for rows.Next() {
		var streamer_id, platform_id string
		if err := rows.Scan(&streamer_id, &platform_id); err != nil {
			return res, utility.NewError(err.Error(), utility.ERR_SQL_DATASCAN)
		}
		res[platform_id] = streamer_id
	}

	return res, nil
}
