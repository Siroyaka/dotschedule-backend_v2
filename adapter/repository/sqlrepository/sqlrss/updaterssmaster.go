package sqlrss

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type UpdateRSSMasterRepository struct {
	updateWrapper sqlwrapper.UpdateWrapper
}

func NewUpdateRSSMasterRepository(sqlHandler abstruct.SqlHandler, query string) UpdateRSSMasterRepository {
	return UpdateRSSMasterRepository{
		updateWrapper: sqlwrapper.NewUpdateWrapper(sqlHandler, query),
	}
}

func (repos UpdateRSSMasterRepository) Execute(data reference.IDWithTime) (reference.DBUpdateResponse, utility.IError) {
	count, id, err := repos.updateWrapper.UpdatePrepare(
		data.Time().ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat()),
		data.Id(),
	)
	if err != nil {
		return reference.DBUpdateResponse{}, err.WrapError("RSS Master Update Error.", "")
	}

	return reference.DBUpdateResponse{
		Count: count,
		Id:    id,
	}, nil
}
