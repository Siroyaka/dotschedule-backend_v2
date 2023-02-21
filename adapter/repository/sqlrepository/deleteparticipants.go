package sqlrepository

import (
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type DeleteParticipantsRepository struct {
	updateWrapper        sqlwrapper.UpdateWrapper
	queryTemplate        string
	replaceTargetsString string
	replacedChar         string
	replacedCharSplitter string
}

func NewDeleteParticipantsRepository(
	sqlHandler abstruct.SqlHandler,
	queryTemplate string,
	replaceTargetsString string,
	replacedChar string,
	replacedCharSplitter string,
) DeleteParticipantsRepository {
	return DeleteParticipantsRepository{
		updateWrapper:        sqlwrapper.NewUpdateWrapper(sqlHandler, ""),
		queryTemplate:        queryTemplate,
		replaceTargetsString: replaceTargetsString,
		replacedChar:         replacedChar,
		replacedCharSplitter: replacedChar,
	}
}

func (repos DeleteParticipantsRepository) Execute(data domain.StreamingParticipants) (reference.DBUpdateResponse, utility.IError) {
	var replacedCharList []string
	for range data.GetList() {
		replacedCharList = append(replacedCharList, repos.replacedChar)
	}
	var replacedString = strings.Join(replacedCharList, repos.replacedCharSplitter)
	queryTemplate := utility.ReplaceConstString(repos.queryTemplate, replacedString, repos.replaceTargetsString)

	repos.updateWrapper.SetQuery(queryTemplate)

	count, id, err := repos.updateWrapper.UpdatePrepare(utility.ToInterfaceSlice(data.StreamingID(), data.Platform(), data.GetList())...)
	if err != nil {
		return reference.DBUpdateResponse{}, err.WrapError()
	}

	return reference.DBUpdateResponse{
		Count: count,
		Id:    id,
	}, nil

}
