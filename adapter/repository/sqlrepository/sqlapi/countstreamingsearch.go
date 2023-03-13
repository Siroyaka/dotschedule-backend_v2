package sqlapi

import (
	"fmt"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference/apireference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type CountStreamingSearchRepository struct {
	selectWrapper        sqlwrapper.SelectWrapper[int]
	mainQuery            string
	subQueryMember       string
	subQueryTags         string
	subQueryTitle        string
	subQueryFrom         string
	subQueryTo           string
	defaultFrom          string
	replaceTargetsString string
	replaceChar          string
	replaceCharSplitter  string
	viewStatus           int
}

func NewCountStreamingSearchRepository(
	sqlHandler abstruct.SqlHandler,
	mainQuery string,
	subQueryMember string,
	subQueryTags string,
	subQueryFrom string,
	subQueryTo string,
	subQueryTitle string,
	defaultFrom string,
	replaceTargetsString string,
	replaceChar string,
	replaceCharSplitter string,
	viewStatus int,
) CountStreamingSearchRepository {
	return CountStreamingSearchRepository{
		selectWrapper:        sqlwrapper.NewSelectWrapper[int](sqlHandler, ""),
		mainQuery:            mainQuery,
		subQueryMember:       subQueryMember,
		subQueryTags:         subQueryTags,
		subQueryFrom:         subQueryFrom,
		subQueryTo:           subQueryTo,
		subQueryTitle:        subQueryTitle,
		defaultFrom:          defaultFrom,
		replaceTargetsString: replaceTargetsString,
		replaceChar:          replaceChar,
		replaceCharSplitter:  replaceCharSplitter,
		viewStatus:           viewStatus,
	}
}

func (repos CountStreamingSearchRepository) scan(s sqlwrapper.IScan) (int, utilerror.IError) {
	var search_length int
	if err := s.Scan(
		&search_length,
	); err != nil {
		return 0, utilerror.New(err.Error(), "")
	}

	return search_length, nil
}

func (repos CountStreamingSearchRepository) createQueryWheres(members []string, from, to wrappedbasics.IWrappedTime, title string) (string, []interface{}) {
	var whereQuerys []string
	var whereValues []interface{}

	whereQuerys = append(whereQuerys, repos.subQueryFrom)

	if from != nil {
		whereValues = append(whereValues, from.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat()))
	} else {
		whereValues = append(whereValues, repos.defaultFrom)
	}

	if to != nil {
		whereQuerys = append(whereQuerys, repos.subQueryTo)
		whereValues = append(whereValues, to.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat()))
	}

	// streaming_id IN (SELECT streaming_id FROM streaming_participants WHERE member_id IN (?, ?) GROUP BY streaming_id HAVING COUNT(streaming_id) = ?)
	if len(members) > 0 {
		var replacedCharList []string
		for _, value := range members {
			if value == "" {
				continue
			}
			replacedCharList = append(replacedCharList, repos.replaceChar)
			whereValues = append(whereValues, value)
		}

		if len(replacedCharList) > 0 {
			var replacedString = strings.Join(replacedCharList, repos.replaceCharSplitter)
			queryTemplate := utility.ReplaceConstString(repos.subQueryMember, replacedString, repos.replaceTargetsString)

			whereValues = append(whereValues, len(members))

			whereQuerys = append(whereQuerys, queryTemplate)
		}
	}

	if len(title) > 0 {
		whereValues = append(whereValues, fmt.Sprintf("%%%s%%", title))
		whereQuerys = append(whereQuerys, repos.subQueryTitle)
	}

	return fmt.Sprintf("WHERE %s", strings.Join(whereQuerys, " AND ")), whereValues
}

func (repos CountStreamingSearchRepository) Execute(data apireference.StreamingSearchValues) (int, utilerror.IError) {
	members, from, to, _, _, title := data.Extract()

	whereQuery, whereValues := repos.createQueryWheres(members, from, to, title)

	queryTemplate := utility.ReplaceConstString(repos.mainQuery, whereQuery, repos.replaceTargetsString)

	repos.selectWrapper.SetQuery(queryTemplate)

	result, err := repos.selectWrapper.SelectPrepare(repos.scan, utility.ToInterfaceSlice(whereValues)...)
	if err != nil {
		return 0, err.WrapError()
	}
	return result[0], nil
}
