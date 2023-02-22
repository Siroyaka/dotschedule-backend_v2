package sqlrepository

import (
	"database/sql"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

type SelectNormalizeSchedulesRepository struct {
	selectWrapper sqlwrapper.SelectWrapper[domain.FullScheduleData]
}

func NewSelectNormalizeSchedulesRepository(sqlHandler abstruct.SqlHandler, query string) SelectNormalizeSchedulesRepository {
	return SelectNormalizeSchedulesRepository{
		selectWrapper: sqlwrapper.NewSelectWrapper[domain.FullScheduleData](sqlHandler, query),
	}
}

func (repos SelectNormalizeSchedulesRepository) scheduleScan(s sqlwrapper.IScan) (domain.FullScheduleData, utilerror.IError) {
	var streaming_id, platform_type, status, insert_at string
	var publish_datetime sql.NullString

	if err := s.Scan(&streaming_id, &platform_type, &status, &publish_datetime, &insert_at); err != nil {
		return domain.NewEmptyFullScheduleData("", ""), utilerror.New(err.Error(), "")
	}

	res := domain.NewEmptyFullScheduleData(streaming_id, platform_type)
	res.Status = status
	res.InsertAt = insert_at
	if !publish_datetime.Valid {
		res.PublishDatetime = publish_datetime.String
	}

	return res, nil
}

func (repos SelectNormalizeSchedulesRepository) Execute(_ reference.VoidStruct) ([]domain.FullScheduleData, utilerror.IError) {
	result, err := repos.selectWrapper.Select(repos.scheduleScan)
	if err != nil {
		return result, err.WrapError()
	}
	return result, nil
}
