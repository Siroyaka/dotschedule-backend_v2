package sqlapi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/domain/apidomain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference/apireference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type SelectStreamingSearchRepository struct {
	selectWrapper        sqlwrapper.SelectWrapper[apireference.ScheduleResponse]
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
	limit                int
}

func NewSelectStreamingSearchRepository(
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
	limit int,
) SelectStreamingSearchRepository {
	return SelectStreamingSearchRepository{
		selectWrapper:        sqlwrapper.NewSelectWrapper[apireference.ScheduleResponse](sqlHandler, ""),
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
		limit:                limit,
	}
}

func (repos SelectStreamingSearchRepository) scan(s sqlwrapper.IScan) (apireference.ScheduleResponse, utilerror.IError) {
	var streaming_id string
	var url string
	var platform string
	var status string
	var publish_datetime string
	var duration int
	var thumbnail string
	var title string
	var description string
	var streamer_name string
	var streamer_id string
	var streamer_icon string
	var streamer_link string
	var participants_data string
	if err := s.Scan(
		&streaming_id,
		&url,
		&platform,
		&status,
		&publish_datetime,
		&duration,
		&thumbnail,
		&title,
		&description,
		&streamer_name,
		&streamer_id,
		&streamer_icon,
		&streamer_link,
		&participants_data,
	); err != nil {
		return apireference.ScheduleResponse{}, utilerror.New(err.Error(), "")
	}

	statusNum, err := strconv.Atoi(status)
	if err != nil {
		logger.Error(err)
		return apireference.ScheduleResponse{}, utilerror.New(err.Error(), "")
	}

	startDate, ierr := wrappedbasics.NewWrappedTimeFromUTC(publish_datetime, wrappedbasics.WrappedTimeProps.DateTimeFormat())
	if ierr != nil {
		return apireference.ScheduleResponse{}, ierr.WrapError()
	}

	streamingData := apidomain.StreamingData{
		ID:          streaming_id,
		URL:         url,
		Platform:    platform,
		Status:      statusNum,
		StartDate:   startDate.ToLocalFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat()),
		Duration:    duration,
		Thumbnail:   thumbnail,
		Title:       title,
		Description: description,
	}

	streamerData := apidomain.StreamerData{
		Name:     streamer_name,
		ID:       streamer_id,
		Icon:     streamer_icon,
		Link:     streamer_link,
		Platform: platform,
	}

	participants, ierr := utility.JsonUnmarshal[[]apidomain.StreamerData](participants_data)
	if ierr != nil {
		participants = []apidomain.StreamerData{}
		fmt.Println("partse error")
		fmt.Println(ierr.Error())
	}

	return apireference.ScheduleResponse{
		StreamingData: streamingData,
		StreamerData:  streamerData,
		Participants:  participants,
	}, nil
}

func (repos SelectStreamingSearchRepository) createQueryWheres(members []string, from, to wrappedbasics.IWrappedTime, title string) (string, []interface{}) {
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

func (repos SelectStreamingSearchRepository) Execute(data apireference.StreamingSearchValues) ([]apireference.ScheduleResponse, utilerror.IError) {
	members, from, to, _, page, title := data.Extract()

	whereQuery, whereValues := repos.createQueryWheres(members, from, to, title)

	offset := 0
	if page > 0 {
		offset = (page - 1) * repos.limit
	}

	queryTemplate := utility.ReplaceConstString(repos.mainQuery, whereQuery, repos.replaceTargetsString)

	repos.selectWrapper.SetQuery(queryTemplate)

	result, err := repos.selectWrapper.SelectPrepare(repos.scan, utility.ToInterfaceSlice(whereValues, repos.limit, offset)...)
	if err != nil {
		return result, err.WrapError()
	}
	return result, nil
}
