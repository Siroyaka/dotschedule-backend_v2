package streamingparticipants

import (
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/streamingparticipants"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type InsertRepository struct {
	sqlHandler           abstruct.SqlHandler
	queryTemplate        string
	replacedTargetString string
	replaceChar          string
	replaceSplitter      string
}

func NewInsertRepository(
	sqlHandler abstruct.SqlHandler,
	queryTemplate, replacedTargetString, replaceChar, replaceSplitter string,
) streamingparticipants.InsertRepository {
	return InsertRepository{
		sqlHandler:           sqlHandler,
		queryTemplate:        queryTemplate,
		replacedTargetString: replacedTargetString,
		replaceChar:          replaceChar,
		replaceSplitter:      replaceSplitter,
	}
}

func (r InsertRepository) insert(streamingId, platform string, list []string) (int64, utility.IError) {
	if len(list) == 0 {
		return 0, nil
	}

	insertAt := wrappedbasics.Now().ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())

	var replaceCharList []string
	for range list {
		replaceCharList = append(replaceCharList, r.replaceChar)
	}

	replaceText := strings.Join(replaceCharList, r.replaceSplitter)

	query := utility.ReplaceConstString(r.queryTemplate, replaceText, r.replacedTargetString)
	stmt, err := r.sqlHandler.Prepare(query)
	if err != nil {
		return 0, utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, query)
	}
	result, err := stmt.Exec(utility.ToInterfaceSlice(streamingId, insertAt, platform, list)...)
	if err != nil {
		return 0, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}
	affectedCount, err := result.RowsAffected()
	if err != nil {
		return affectedCount, utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
	}

	return affectedCount, nil

}

func (r InsertRepository) InsertList(streamingId, platformType string, memberIdList ...string) (int64, utility.IError) {
	result, err := r.insert(streamingId, platformType, memberIdList)
	if err != nil {
		return result, err.WrapError()
	}
	return result, nil
}

func (r InsertRepository) InsertStreamingParticipants(data domain.StreamingParticipants) (int64, utility.IError) {
	if data.IsEmpty() {
		return 0, nil
	}
	result, err := r.insert(data.StreamingID(), data.Platform(), data.GetList())
	if err != nil {
		return result, err.WrapError()
	}
	return result, nil
}
