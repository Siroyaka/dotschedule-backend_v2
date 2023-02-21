package sqlrepository

import (
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type InsertParticipantsRepository struct {
	updateWrapper        sqlwrapper.UpdateWrapper
	queryTemplate        string
	replacedTargetString string
	replaceChar          string
	replaceSplitter      string
}

func NewInsertParticipantsRepository(
	sqlHandler abstruct.SqlHandler,
	queryTemplate, replacedTargetString, replaceChar, replaceSplitter string,
) InsertParticipantsRepository {
	return InsertParticipantsRepository{
		updateWrapper:        sqlwrapper.NewUpdateWrapper(sqlHandler, ""),
		queryTemplate:        queryTemplate,
		replacedTargetString: replacedTargetString,
		replaceChar:          replaceChar,
		replaceSplitter:      replaceSplitter,
	}
}

func (repos InsertParticipantsRepository) Execute(data domain.StreamingParticipants) (reference.DBUpdateResponse, utility.IError) {
	insertAt := wrappedbasics.Now().ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())

	var replaceCharList []string
	for range data.GetList() {
		replaceCharList = append(replaceCharList, repos.replaceChar)
	}

	replaceText := strings.Join(replaceCharList, repos.replaceSplitter)

	query := utility.ReplaceConstString(repos.queryTemplate, replaceText, repos.replacedTargetString)

	repos.updateWrapper.SetQuery(query)

	count, id, err := repos.updateWrapper.UpdatePrepare(utility.ToInterfaceSlice(data.StreamingID(), insertAt, data.Platform(), data.GetList())...)
	if err != nil {
		return reference.DBUpdateResponse{}, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	return reference.DBUpdateResponse{
		Count: count,
		Id:    id,
	}, nil
}
