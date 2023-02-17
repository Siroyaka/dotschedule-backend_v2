package sqlrepository

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type InsertRSSSeedScheduleRepository struct {
	updateWrapper sqlwrapper.UpdateWrapper
}

func NewInsertRSSSeedScheduleRepository(sqlHandler abstruct.SqlHandler, query string) InsertRSSSeedScheduleRepository {
	return InsertRSSSeedScheduleRepository{
		updateWrapper: sqlwrapper.NewUpdateWrapper(sqlHandler, query),
	}
}

func (repos InsertRSSSeedScheduleRepository) Execute(data domain.SeedSchedule) (reference.DBUpdateResponse, utility.IError) {
	now := wrappedbasics.Now(wrappedbasics.WrappedTimeProps.LocalLocation())

	nowString := now.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())

	count, id, err := repos.updateWrapper.UpdatePrepare(data.GetID(), data.GetPlatformType(), data.GetStatus(), nowString, nowString, data.GetVisibleStatus(), data.GetCompleteStatus())

	if err != nil {
		err = err.WrapError()
	}

	res := reference.DBUpdateResponse{
		Count: count,
		Id:    id,
	}
	return res, err
}
