package migration

import (
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	migration "github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/migration"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type GetRepository struct {
	sqlHandler                  abstruct.SqlHandler
	checkQueryTemplate          string
	replacedChar                string
	replacedCharSplitter        string
	replaceTargetsString        string
	streamerMasterQueryTemplate string
}

func NewGetRepository(
	sqlHandler abstruct.SqlHandler,
	checkQueryTemplate,
	replacedChar,
	replacedCharSplitter,
	replaceTargetsString,
	streamerMasterQueryTemplate string,
) migration.GetDataRepository {
	return GetRepository{
		sqlHandler:                  sqlHandler,
		checkQueryTemplate:          checkQueryTemplate,
		replacedChar:                replacedChar,
		replacedCharSplitter:        replacedCharSplitter,
		replaceTargetsString:        replaceTargetsString,
		streamerMasterQueryTemplate: streamerMasterQueryTemplate,
	}
}

func (repos GetRepository) GetStreamerMaster(platform string) (result map[string]domain.GroupStreamerData, ierr utility.IError) {
	result = make(map[string]domain.GroupStreamerData)
	sqmt, err := repos.sqlHandler.Prepare(repos.streamerMasterQueryTemplate)
	if err != nil {
		ierr = utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, repos.streamerMasterQueryTemplate)
		return
	}
	defer sqmt.Close()

	row, err := sqmt.Query(platform)
	if err != nil {
		ierr = utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
		return
	}
	defer row.Close()

	for row.Next() {
		var streamer_id string
		var platform_id string
		var streamer_name string
		row.Scan(&streamer_id, &platform_id, &streamer_name)
		result[platform_id] = domain.GroupStreamerData{
			StreamerID:   streamer_id,
			StreamerName: streamer_name,
		}
	}

	ierr = nil
	return
}

func (repos GetRepository) GetIDSet(list []string, platform string) (result utility.HashSet[string], ierr utility.IError) {
	result = utility.NewHashSet[string]()
	var replacedCharList []string
	for range list {
		replacedCharList = append(replacedCharList, repos.replacedChar)
	}
	var replacedString = strings.Join(replacedCharList, repos.replacedCharSplitter)
	queryTemplate := utility.ReplaceConstString(repos.checkQueryTemplate, replacedString, repos.replaceTargetsString)

	sqmt, err := repos.sqlHandler.Prepare(queryTemplate)
	if err != nil {
		ierr = utility.NewError(err.Error(), utility.ERR_SQL_PREPARE, queryTemplate)
		return
	}
	defer sqmt.Close()

	row, err := sqmt.Query(utility.ToInterfaceSlice(platform, list)...)
	if err != nil {
		ierr = utility.NewError(err.Error(), utility.ERR_SQL_QUERY)
		return
	}
	defer row.Close()
	for row.Next() {
		var id string
		row.Scan(&id)
		result.Set(id)
	}
	ierr = nil
	return
}
