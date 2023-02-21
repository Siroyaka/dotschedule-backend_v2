package sqlrepository

import (
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type UpdateScheduleOnlyCompleteTo0Repository struct {
	updateWrapper        sqlwrapper.UpdateWrapper
	updateQueryTemplate  string
	replaceTargetsString string
	replaceChar          string
	replaceCharSplitter  string
}

func NewUpdateScheduleOnlyCompleteTo0Repository(
	sqlHandler abstruct.SqlHandler,
	updateQueryTemplate, replaceTargetsString, replaceChar, replaceCharSplitter string,
) UpdateScheduleOnlyCompleteTo0Repository {
	return UpdateScheduleOnlyCompleteTo0Repository{
		updateWrapper:        sqlwrapper.NewUpdateWrapper(sqlHandler, ""),
		updateQueryTemplate:  updateQueryTemplate,
		replaceTargetsString: replaceTargetsString,
		replaceChar:          replaceChar,
		replaceCharSplitter:  replaceCharSplitter,
	}
}

func (repos UpdateScheduleOnlyCompleteTo0Repository) Execute(data reference.StreamingIDListWithPlatformID) (reference.DBUpdateResponse, utility.IError) {
	completeStatus := 0
	idList, platformType := data.Extract()
	if len(idList) == 0 {
		return reference.DBUpdateResponse{}, nil
	}

	now := wrappedbasics.Now()
	var replacedCharList []string
	for i := 0; i < len(idList); i++ {
		replacedCharList = append(replacedCharList, repos.replaceChar)
	}
	var replacedString = strings.Join(replacedCharList, repos.replaceCharSplitter)
	queryTemplate := utility.ReplaceConstString(repos.updateQueryTemplate, replacedString, repos.replaceTargetsString)
	repos.updateWrapper.SetQuery(queryTemplate)

	count, id, err := repos.updateWrapper.UpdatePrepare(utility.ToInterfaceSlice(completeStatus, now.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat()), platformType, idList)...)

	if err != nil {
		return reference.DBUpdateResponse{Count: count, Id: id}, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	return reference.DBUpdateResponse{Count: count, Id: id}, nil
}
