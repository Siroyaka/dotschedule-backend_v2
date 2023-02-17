package sqlrss

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type SelectRSSMasterRepository struct {
	selectWrapper sqlwrapper.SelectWrapper[domain.RSSMaster]
}

func NewSelectRSSMasterRepository(sqlHandler abstruct.SqlHandler, query string) SelectRSSMasterRepository {
	return SelectRSSMasterRepository{
		selectWrapper: sqlwrapper.NewSelectWrapper[domain.RSSMaster](sqlHandler, query),
	}
}

func (repos SelectRSSMasterRepository) scan(s sqlwrapper.IScan) (domain.RSSMaster, utility.IError) {
	var streamer_id string
	var streamer_name string
	var rss_url string
	var publish_datetime string
	if err := s.Scan(&streamer_id, &streamer_name, &rss_url, &publish_datetime); err != nil {
		return domain.RSSMaster{}, utility.NewError(err.Error(), "")
	}

	dt, err := wrappedbasics.NewWrappedTimeFromUTC(publish_datetime, wrappedbasics.WrappedTimeProps.DateTimeFormat())
	if err != nil {
		return domain.RSSMaster{}, err.WrapError()
	}

	return domain.NewRSSMaster(
		streamer_id,
		streamer_name,
		rss_url,
		dt,
	), nil
}

func (repos SelectRSSMasterRepository) Execute(_ reference.VoidStruct) ([]domain.RSSMaster, utility.IError) {
	result, err := repos.selectWrapper.Select(repos.scan)
	if err != nil {
		return result, err.WrapError()
	}
	return result, nil
}
