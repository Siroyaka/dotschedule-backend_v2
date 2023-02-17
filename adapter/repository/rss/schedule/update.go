package rssschedule

import (
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	rssschedule "github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/rss/schedule"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type UpdateRepository struct {
	sqlHandler           abstruct.SqlHandler
	updateQueryTemplate  string
	replaceTargetsString string
	replaceChar          string
	replaceCharSplitter  string
}

func NewUpdateRepository(
	sqlHandler abstruct.SqlHandler,
	updateQueryTemplate string,
	replaceTargetsString, replaceChar, replaceCharSplitter string,
) rssschedule.UpdateRepository {
	return UpdateRepository{
		sqlHandler:           sqlHandler,
		updateQueryTemplate:  updateQueryTemplate,
		replaceTargetsString: replaceTargetsString,
		replaceChar:          replaceChar,
		replaceCharSplitter:  replaceCharSplitter,
	}
}

func (r UpdateRepository) Update(idList []string, platformType string, completeStatus int) utility.IError {
	if len(idList) == 0 {
		return nil
	}

	now := wrappedbasics.Now(wrappedbasics.WrappedTimeProps.LocalLocation())
	var replacedCharList []string
	for i := 0; i < len(idList); i++ {
		replacedCharList = append(replacedCharList, r.replaceChar)
	}
	var replacedString = strings.Join(replacedCharList, r.replaceCharSplitter)
	queryTemplate := utility.ReplaceConstString(r.updateQueryTemplate, replacedString, r.replaceTargetsString)
	sqmt, err := r.sqlHandler.Prepare(queryTemplate)
	if err != nil {
		return utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, queryTemplate)
	}
	defer sqmt.Close()
	_, err = sqmt.Exec(utility.ToInterfaceSlice(completeStatus, now.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat()), platformType, idList)...)
	if err != nil {
		return utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	return nil
}
