package sqlrepository

import (
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

type SelectIDsFromIDsRepository struct {
	selectWrapper        sqlwrapper.SelectWrapper[string]
	checkQueryTemplate   string
	replaceTargetsString string
	replaceChar          string
	replaceCharSplitter  string
}

func NewSelectIDsFromIDsRepository(
	sqlHandler abstruct.SqlHandler,
	checkQueryTemplate, replaceTargetsString, replaceChar, replaceCharSplitter string,
) SelectIDsFromIDsRepository {
	return SelectIDsFromIDsRepository{
		selectWrapper:        sqlwrapper.NewSelectWrapper[string](sqlHandler, utility.EmptyString),
		checkQueryTemplate:   checkQueryTemplate,
		replaceTargetsString: replaceTargetsString,
		replaceChar:          replaceChar,
		replaceCharSplitter:  replaceCharSplitter,
	}
}

func (repos SelectIDsFromIDsRepository) scanable(s sqlwrapper.IScan) (string, utilerror.IError) {
	var id string
	var is_complete_data int
	if err := s.Scan(&id, &is_complete_data); err != nil {
		return utility.EmptyString, utilerror.New(err.Error(), "")
	}
	return id, nil
}

func (repos SelectIDsFromIDsRepository) Execute(data reference.StreamingIDListWithPlatformID) (utility.HashSet[string], utilerror.IError) {
	idList, platform := data.Extract()

	var replacedCharList []string
	for i := 0; i < len(idList); i++ {
		replacedCharList = append(replacedCharList, repos.replaceChar)
	}

	var replacedString = strings.Join(replacedCharList, repos.replaceCharSplitter)
	queryTemplate := utility.ReplaceConstString(repos.checkQueryTemplate, replacedString, repos.replaceTargetsString)

	repos.selectWrapper.SetQuery(queryTemplate)

	hashset := utility.NewHashSet[string]()

	list, err := repos.selectWrapper.SelectPrepare(repos.scanable, utility.ToInterfaceSlice(platform, idList)...)
	if err != nil {
		return hashset, err.WrapError()
	}

	for _, id := range list {
		hashset.Set(id)
	}

	return hashset, nil
}
