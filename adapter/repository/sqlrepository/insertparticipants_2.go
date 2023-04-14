package sqlrepository

import (
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type InsertParticipantsRepository2 struct {
	updateWrapper        sqlwrapper.UpdateWrapper
	queryTemplate        string
	replacedTargetString string
	replaceChar          string
	replaceSplitter      string
}

func NewInsertParticipantsRepository2(
	sqlHandler abstruct.SqlHandler,
	queryTemplate, replacedTargetString, replaceChar, replaceSplitter string,
) InsertParticipantsRepository2 {
	return InsertParticipantsRepository2{
		updateWrapper:        sqlwrapper.NewUpdateWrapper(sqlHandler, ""),
		queryTemplate:        queryTemplate,
		replacedTargetString: replacedTargetString,
		replaceChar:          replaceChar,
		replaceSplitter:      replaceSplitter,
	}
}

func (repos InsertParticipantsRepository2) Execute(data domain.StreamingParticipants) (reference.DBUpdateResponse, utilerror.IError) {
	insertAt := wrappedbasics.Now().ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())

	var replaceCharList []string
	for range data.GetList() {
		replaceCharList = append(replaceCharList, repos.replaceChar)
	}

	replaceText := strings.Join(replaceCharList, repos.replaceSplitter)

	query := utility.ReplaceConstString(repos.queryTemplate, replaceText, repos.replacedTargetString)

	repos.updateWrapper.SetQuery(query)

	count, id, err := repos.updateWrapper.UpdatePrepare(utility.ToInterfaceSlice(data.StreamingID(), data.Platform(), insertAt, data.GetList())...)
	if err != nil {
		return reference.DBUpdateResponse{}, utilerror.New(err.Error(), utilerror.ERR_SQL_QUERY)
	}

	return reference.DBUpdateResponse{
		Count: count,
		Id:    id,
	}, nil
}
