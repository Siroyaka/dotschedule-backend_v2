package sqlrepository

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/dbmodels"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

type SelectAScheduleParticipantsRepository struct {
	selectWrapper sqlwrapper.SelectWrapper[dbmodels.KeyValue[string, string]]
}

func NewSelectAScheduleParticipantsRepository(sqlHandler abstruct.SqlHandler, query string) SelectAScheduleParticipantsRepository {
	return SelectAScheduleParticipantsRepository{
		selectWrapper: sqlwrapper.NewSelectWrapper[dbmodels.KeyValue[string, string]](sqlHandler, query),
	}
}

func (repos SelectAScheduleParticipantsRepository) participantsIdNameFromDb(s sqlwrapper.IScan) (dbmodels.KeyValue[string, string], utilerror.IError) {
	var id, name string
	if err := s.Scan(&id, &name); err != nil {
		return dbmodels.EmptyKeyValue[string, string](), utilerror.New(err.Error(), "")
	}
	return dbmodels.NewKeyValue(id, name), nil
}

func (repos SelectAScheduleParticipantsRepository) Execute(data reference.StreamingIDWithPlatformType) ([]dbmodels.KeyValue[string, string], utilerror.IError) {
	streamingId, platformType := data.Extract()

	result, err := repos.selectWrapper.SelectPrepare(repos.participantsIdNameFromDb, streamingId, platformType)

	if err != nil {
		return result, err.WrapError()
	}

	return result, nil
}
