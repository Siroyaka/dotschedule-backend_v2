package rssschedule

import (
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	rssschedule "github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/rss/schedule"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetRepository struct {
	sqlHandler           abstruct.SqlHandler
	common               utility.Common
	checkQueryTemplate   string
	replaceTargetsString string
	replaceChar          string
	replaceCharSplitter  string
}

func NewGetRepository(
	sqlHandler abstruct.SqlHandler,
	common utility.Common,
	checkQueryTemplate,
	replaceTargetsString, replaceChar, replaceCharSplitter string,
) rssschedule.GetRepository {
	return GetRepository{
		sqlHandler:           sqlHandler,
		common:               common,
		checkQueryTemplate:   checkQueryTemplate,
		replaceTargetsString: replaceTargetsString,
		replaceChar:          replaceChar,
		replaceCharSplitter:  replaceCharSplitter,
	}
}

func test(s string, i int) (string, bool) {
	return "", true
}

func (r GetRepository) Get(idList []string, platform string, function rssschedule.ScheduleIsCompleteConverter) (domain.ScheduleIDWithIsComplete, utility.IError) {
	scheduleIsComplete := domain.NewScheduleIDWithIsComplete()
	var replacedCharList []string
	for i := 0; i < len(idList); i++ {
		replacedCharList = append(replacedCharList, r.replaceChar)
	}
	var replacedString = strings.Join(replacedCharList, r.replaceCharSplitter)
	queryTemplate := utility.ReplaceConstString(r.checkQueryTemplate, replacedString, r.replaceTargetsString)
	sqmt, err := r.sqlHandler.Prepare(queryTemplate)
	if err != nil {
		return scheduleIsComplete, utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, queryTemplate)
	}
	defer sqmt.Close()

	row, err := sqmt.Query(utility.ToInterfaceSlice(platform, idList)...)
	if err != nil {
		return scheduleIsComplete, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}
	defer row.Close()

	for row.Next() {
		var id string
		var is_complete_data int
		row.Scan(&id, &is_complete_data)
		scheduleIsComplete.Add(function(id, is_complete_data))
	}
	return scheduleIsComplete, nil
}
