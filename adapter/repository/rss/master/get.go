package rssmaster

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	rssmaster "github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/rss/master"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetRepository struct {
	sqlHandler     abstruct.SqlHandler
	rssMasterQuery string
}

func NewGetRepository(handler abstruct.SqlHandler, rssMasterQuery string) rssmaster.GetRepository {
	return GetRepository{
		sqlHandler:     handler,
		rssMasterQuery: rssMasterQuery,
	}
}

func (r GetRepository) Get(adapter rssmaster.MasterDataAdapter) ([]domain.RSSMaster, utility.IError) {
	var list []domain.RSSMaster

	rows, err := r.sqlHandler.Query(r.rssMasterQuery)
	if err != nil {
		return list, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}
	defer rows.Close()

	for rows.Next() {
		var streamer_id string
		var streamer_name string
		var rss_url string
		var publish_datetime string
		err := rows.Scan(&streamer_id, &streamer_name, &rss_url, &publish_datetime)
		if err != nil {
			utility.LogError(utility.NewError(err.Error(), utility.ERR_SQL_DATASCAN))
			continue
		}
		data, u_err := adapter(streamer_id, streamer_name, rss_url, publish_datetime)
		if u_err != nil {
			utility.LogError(u_err.WrapError())
			continue
		}
		list = append(list, data)
	}

	return list, nil
}
