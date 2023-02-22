package sqlrepository

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

type SelectStreamerMasterWithPlatformMasterRepository struct {
	selectWrapper sqlwrapper.SelectWrapper[domain.StreamerMasterWithPlatformData]
}

func NewSelectStreamerMasterWithPlatformMasterRepository(sqlHandler abstruct.SqlHandler, query string) SelectStreamerMasterWithPlatformMasterRepository {
	return SelectStreamerMasterWithPlatformMasterRepository{
		selectWrapper: sqlwrapper.NewSelectWrapper[domain.StreamerMasterWithPlatformData](sqlHandler, query),
	}
}

func (repos SelectStreamerMasterWithPlatformMasterRepository) makeScan(platformType string) func(sqlwrapper.IScan) (domain.StreamerMasterWithPlatformData, utilerror.IError) {
	f := func(s sqlwrapper.IScan) (domain.StreamerMasterWithPlatformData, utilerror.IError) {
		var streamer_id, platform_id, streamer_name string

		if err := s.Scan(&streamer_id, &platform_id, &streamer_name); err != nil {
			return domain.NewStreamerMasterWithPlatformData(""), utilerror.New(err.Error(), "")
		}

		res := domain.NewStreamerMasterWithPlatformData(streamer_id)
		res.StreamerName = streamer_name
		res.PlatformData[platformType] = domain.StreamerPlatformMaster{
			StreamerID:   streamer_id,
			PlatformID:   platform_id,
			PlatformType: platformType,
		}
		return res, nil
	}
	return f
}

func (repos SelectStreamerMasterWithPlatformMasterRepository) Execute(platformType string) ([]domain.StreamerMasterWithPlatformData, utilerror.IError) {
	result, err := repos.selectWrapper.SelectPrepare(repos.makeScan(platformType), platformType)
	if err != nil {
		return result, err.WrapError()
	}
	return result, nil
}
