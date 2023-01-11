package usecase

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/fileio"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/fullschedule"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/migration"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/streamingparticipants"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

const (
	diStartIndex = -1
)

type DataMigrationInteractor struct {
	fileReaderRepos         fileio.ReaderRepository[domain.JsonScheduleList]
	getMigrationDataRepos   migration.GetDataRepository
	insertFullScheduleRepos fullschedule.InsertRepository
	common                  utility.Common
	dataFilePathList        []string
	index                   int
	insertParticipantsRepos streamingparticipants.InsertRepository
	platform                string
}

func NewDataMigrationInteractor(
	common utility.Common,
	fRepos fileio.ReaderRepository[domain.JsonScheduleList],
	dataFileDirectoryPath, platform string,
	gRepos migration.GetDataRepository,
	ifRepos fullschedule.InsertRepository,
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
		insertFullScheduleRepos: ifRepos,
		common:                  common,
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
	common utility.Common,
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
		itemDate, err := common.CreateNewWrappedTimeFromLocal(item.StartDate)
		if err != nil {
			utility.LogError(err.WrapError(fmt.Sprintf("id: %s", item.VideoID)))

			// if datetime parse error then status = -1 because it must be updated after migration.
			videoStatus = -1
		} else {
			publishDatetime = itemDate.ToUTCFormatString()
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
	insertDataList := pickupInsertTarget(data, master, registered, intr.common, intr.platform)

	// data insert
	for _, insertData := range insertDataList {
		now, err := intr.common.Now()
		if err != nil {
			utility.LogFatal(err.WrapError(fmt.Sprintf("irregular error. streamingID: %s", insertData.StreamingID)))
			errCount++
			continue
		}
		cnt, err := intr.insertFullScheduleRepos.Insert(insertData.FullScheduleData, now)
		if err != nil {
			utility.LogError(err.WrapError(fmt.Sprintf("schedule insert error. streamingID: %s", insertData.StreamingID)))
			errCount++
			continue
		}
		if cnt == 0 {
			utility.LogError(utility.NewError(fmt.Sprintf("insert data count = 0. streamingID: %s", insertData.StreamingID), ""))
			errCount++
			continue
		}
		utility.LogInfo(fmt.Sprintf("schedule insert finished. streamingID: %s", insertData.StreamingID))

		affectedCount, err := intr.insertParticipantsRepos.InsertList(insertData.StreamingID, insertData.PlatformType, now, insertData.PlatformIdList...)

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
