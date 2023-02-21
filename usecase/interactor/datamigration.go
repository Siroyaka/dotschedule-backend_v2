package interactor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/fileio"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/migration"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/streamingparticipants"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/reference"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

const (
	diStartIndex = -1
)

type DataMigrationInteractor struct {
	fileReaderRepos         fileio.ReaderRepository[domain.JsonScheduleList]
	getMigrationDataRepos   migration.GetDataRepository
	insertFullScheduleRepos abstruct.RepositoryRequest[domain.FullScheduleData, reference.DBUpdateResponse]
	dataFilePathList        []string
	index                   int
	insertParticipantsRepos streamingparticipants.InsertRepository
	platform                string
}

func NewDataMigrationInteractor(
	fRepos fileio.ReaderRepository[domain.JsonScheduleList],
	dataFileDirectoryPath, platform string,
	gRepos migration.GetDataRepository,
	insertFullScheduleRepos abstruct.RepositoryRequest[domain.FullScheduleData, reference.DBUpdateResponse],
	iaRepos streamingparticipants.InsertRepository,
) DataMigrationInteractor {
	absPath, err := filepath.Abs(dataFileDirectoryPath)
	if err != nil {
		panic(err)
	}
	items, err := fRepos.FileList(absPath)
	if err != nil {
		panic(err)
	}
	var dataFilePathList []string
	for _, v := range items {
		dataFilePathList = append(dataFilePathList, filepath.Join(absPath, v))
	}
	return DataMigrationInteractor{
		fileReaderRepos:         fRepos,
		getMigrationDataRepos:   gRepos,
		insertFullScheduleRepos: insertFullScheduleRepos,
		dataFilePathList:        dataFilePathList,
		index:                   diStartIndex,
		insertParticipantsRepos: iaRepos,
		platform:                platform,
	}
}

func (intr DataMigrationInteractor) Len() int {
	return len(intr.dataFilePathList)
}

func (intr *DataMigrationInteractor) Next() bool {
	intr.index++
	return intr.index < intr.Len()
}

func pickupInsertTarget(
	data domain.JsonScheduleList,
	master map[string]domain.GroupStreamerData,
	registered utility.HashSet[string],
	platformName string,
) (list []domain.FullScheduleWithPlatformParticipantsData) {
	for _, item := range data {
		if registered.Has(item.VideoID) {
			utility.LogDebug(fmt.Sprintf("registered data: %s", item.VideoID))
			continue
		}

		var streamerName string
		var streamerId string
		if v, ok := master[item.StreamerID]; ok {
			streamerName = v.StreamerName
			streamerId = v.StreamerID
		} else {
			streamerName = item.StreamerName
		}

		publishDatetime := ""
		videoStatus := item.VideoStatus
		itemDate, err := wrappedbasics.NewWrappedTimeFromLocal(item.StartDate, wrappedbasics.WrappedTimeProps.DateTimeFormat())
		if err != nil {
			utility.LogError(err.WrapError(fmt.Sprintf("id: %s", item.VideoID)))

			// if datetime parse error then status = -1 because it must be updated after migration.
			videoStatus = -1
		} else {
			publishDatetime = itemDate.ToUTCFormatString(wrappedbasics.WrappedTimeProps.DateTimeFormat())
		}

		var participants []string
		for key := range item.Participants {
			participants = append(participants, key)
		}

		data := domain.NewEmptyFullScheduleWithPlatformParticipantsData(item.VideoID, platformName).ImportStatusFromFirestore(videoStatus)
		data.Url = item.VideoLink
		data.StreamerName = streamerName
		data.StreamerID = streamerId
		data.Title = item.VideoTitle
		data.PublishDatetime = publishDatetime
		data.Duration = item.Duration
		data.ThumbnailLink = item.Thumbnail
		data = data.AppendParticipants(participants...)

		list = append(list, data)
	}
	return
}

func (intr DataMigrationInteractor) Migration() (totalCount, insertCount, registeredCount, errCount int, ierr utility.IError) {
	if intr.index >= intr.Len() {
		ierr = utility.NewError("", utility.ERR_OUTOFINDEX)
		return
	}

	utility.LogInfo(fmt.Sprintf("data migration start: %s", intr.dataFilePathList[intr.index]))
	// file data import
	data, err := intr.fileReaderRepos.ReadJson(intr.dataFilePathList[intr.index])
	if err != nil {
		ierr = err.WrapError()
		return
	}
	totalCount = len(data)

	// get streamerMaster
	master, err := intr.getMigrationDataRepos.GetStreamerMaster(intr.platform)
	if err != nil {
		ierr = err.WrapError()
		return
	}

	// file data id check and choise data
	idList := []string{}
	for _, v := range data {
		idList = append(idList, v.VideoID)
	}
	registered, err := intr.getMigrationDataRepos.GetIDSet(idList, intr.platform)
	registeredCount = len(registered.List())
	if err != nil {
		ierr = err.WrapError()
		return
	}
	insertDataList := pickupInsertTarget(data, master, registered, intr.platform)

	// data insert
	for _, insertData := range insertDataList {

		if result, err := intr.insertFullScheduleRepos.Execute(insertData.FullScheduleData); err != nil {
			utility.LogError(err.WrapError(fmt.Sprintf("schedule insert error. streamingID: %s", insertData.StreamingID)))
			errCount++
			continue
		} else if result.Count == 0 {
			utility.LogError(utility.NewError(fmt.Sprintf("insert data count = 0. streamingID: %s", insertData.StreamingID), ""))
			errCount++
			continue
		}

		utility.LogInfo(fmt.Sprintf("schedule insert finished. streamingID: %s", insertData.StreamingID))

		affectedCount, err := intr.insertParticipantsRepos.InsertList(insertData.StreamingID, insertData.PlatformType, insertData.PlatformIdList...)

		if err != nil {
			utility.LogError(err.WrapError(fmt.Sprintf("participants data insert error. streamingID: %s", insertData.StreamingID)))
			errCount++
			continue
		}

		if affectedCount != int64(len(insertData.PlatformIdList)) {
			utility.LogInfo(fmt.Sprintf("participants data insert count wrong. result count: %d, list count: %d, id: %s", affectedCount, len(insertData.PlatformIdList), insertData.StreamingID))
			errCount++
			continue
		}
		utility.LogInfo(fmt.Sprintf("participants data insert finished. streamingId: %s, count: %d, id: [%s]", insertData.StreamingID, affectedCount, strings.Join(insertData.PlatformIdList, ",")))
		insertCount++
	}
	utility.LogInfo(fmt.Sprintf("data migration end: %s, total: %d, insert: %d, registered: %d, error: %d", intr.dataFilePathList[intr.index], totalCount, insertCount, registeredCount, errCount))
	ierr = nil
	return
}
