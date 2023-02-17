package sqlrepository

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference/participants"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type InsertSingleParticipantsRepository struct {
	updateWrapper sqlwrapper.UpdateWrapper
}

func NewInsertSingleParticipantsRepository(sqlHandler abstruct.SqlHandler, query string) InsertSingleParticipantsRepository {
	return InsertSingleParticipantsRepository{
		updateWrapper: sqlwrapper.NewUpdateWrapper(sqlHandler, query),
	}
}

func (repos InsertSingleParticipantsRepository) Execute(data participants.SingleInsertData) (reference.DBUpdateResponse, utility.IError) {
	streamingId, platformType, streamerId, updateAt := data.Extract()
	count, id, err := repos.updateWrapper.UpdatePrepare(streamingId, platformType, streamerId, updateAt.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat()))

	if err != nil {
		err = err.WrapError()
	}

	res := reference.DBUpdateResponse{
		Count: count,
		Id:    id,
	}
	return res, err
}
