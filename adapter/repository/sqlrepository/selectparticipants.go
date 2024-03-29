package sqlrepository

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

type SelectParticipantsRepository struct {
	selectwrapper sqlwrapper.SelectWrapper[string]
}

func NewSelectParticipantsRepository(sqlHandler abstruct.SqlHandler, queryTemplate string) SelectParticipantsRepository {
	return SelectParticipantsRepository{
		selectwrapper: sqlwrapper.NewSelectWrapper[string](sqlHandler, queryTemplate),
	}
}

func (repos SelectParticipantsRepository) scan(s sqlwrapper.IScan) (string, utilerror.IError) {
	var member_id string
	err := s.Scan(&member_id)
	if err != nil {
		return utility.EmptyString, utilerror.New(err.Error(), utilerror.ERR_SQL_DATASCAN)
	}
	return member_id, nil
}

func (repos SelectParticipantsRepository) Execute(data reference.StreamingIDWithPlatformType) (domain.StreamingParticipants, utilerror.IError) {
	streamingId, platform := data.Extract()

	list, err := repos.selectwrapper.SelectPrepare(repos.scan, streamingId, platform)
	if err != nil {
		return domain.EmptyStreamingParticipants(), utilerror.New(err.Error(), utilerror.ERR_SQL_QUERY)
	}

	result := domain.NewStreamingParticipants(streamingId, platform)
	for _, id := range list {
		result = result.Add(id)
	}
	return result, nil
}
