package streamingparticipants

import (
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/streamingparticipants"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type DeleteRepository struct {
	sqlHandler           abstruct.SqlHandler
	queryTemplate        string
	replaceTargetsString string
	replacedChar         string
	replacedCharSplitter string
}

func NewDeleteStreamingParticipants(
	sqlHandler abstruct.SqlHandler,
	queryTemplate string,
	replaceTargetsString string,
	replacedChar string,
	replacedCharSplitter string,
) streamingparticipants.DeleteRepository {
	return DeleteRepository{
		sqlHandler:           sqlHandler,
		queryTemplate:        queryTemplate,
		replaceTargetsString: replaceTargetsString,
		replacedChar:         replacedChar,
		replacedCharSplitter: replacedChar,
	}
}

func (r DeleteRepository) Delete(data domain.StreamingParticipants) (int64, utility.IError) {
	var replacedCharList []string
	for range data.GetList() {
		replacedCharList = append(replacedCharList, r.replacedChar)
	}
	var replacedString = strings.Join(replacedCharList, r.replacedCharSplitter)
	queryTemplate := utility.ReplaceConstString(r.queryTemplate, replacedString, r.replaceTargetsString)

	sqmt, err := r.sqlHandler.Prepare(queryTemplate)
	if err != nil {
		return 0, utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, r.queryTemplate)
	}
	defer sqmt.Close()

	res, err := sqmt.Exec(utility.ToInterfaceSlice(data.StreamingID(), data.Platform(), data.GetList())...)
	if err != nil {
		return 0, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return cnt, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	return cnt, nil
}
